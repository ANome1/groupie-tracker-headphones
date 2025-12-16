package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const DeezerAPIURL = "https://api.deezer.com"

// Cache for playlists to avoid hitting API every round
var (
	playlistCache = make(map[string][]DeezerTrack)
	cacheMutex    sync.RWMutex
)

// DeezerTrack represents a track from Deezer API
type DeezerTrack struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Preview string `json:"preview"`
	Artist  struct {
		Name string `json:"name"`
	} `json:"artist"`
	Album struct {
		CoverMedium string `json:"cover_medium"`
	} `json:"album"`
}

type DeezerPlaylistResponse struct {
	Data  []DeezerTrack `json:"data"`
	Total int           `json:"total"`
}

// GetRandomTrackFromDeezerPlaylist fetches a random track from a Deezer playlist
func GetRandomTrackFromDeezerPlaylist(playlistID string) (*DeezerTrack, error) {
	// Check cache first
	cacheMutex.RLock()
	tracks, found := playlistCache[playlistID]
	cacheMutex.RUnlock()

	if found && len(tracks) > 0 {
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(tracks))
		return &tracks[randomIndex], nil
	}

	// Not in cache, fetch from API
	var lastErr error
	for i := 0; i < 3; i++ {
		tracks, err := fetchTracksFromAPI(playlistID)
		if err == nil {
			// Update cache
			cacheMutex.Lock()
			playlistCache[playlistID] = tracks
			cacheMutex.Unlock()

			rand.Seed(time.Now().UnixNano())
			randomIndex := rand.Intn(len(tracks))
			return &tracks[randomIndex], nil
		}
		lastErr = err
		time.Sleep(500 * time.Millisecond)
	}

	return nil, fmt.Errorf("failed to fetch tracks after retries: %v", lastErr)
}

func fetchTracksFromAPI(playlistID string) ([]DeezerTrack, error) {
	url := fmt.Sprintf("%s/playlist/%s/tracks?limit=100", DeezerAPIURL, playlistID)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlist: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var playlist DeezerPlaylistResponse
	if err := json.Unmarshal(body, &playlist); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if len(playlist.Data) == 0 {
		return nil, fmt.Errorf("playlist is empty")
	}

	return playlist.Data, nil
}

// ClearCache clears the playlist cache
func ClearDeezerCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	playlistCache = make(map[string][]DeezerTrack)
}
