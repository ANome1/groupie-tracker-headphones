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

var (
	playlistCache = make(map[string][]DeezerTrack)
	cacheMutex    sync.RWMutex
)

var PlaylistMappings = map[string]string{
	"116": "1677006641",
	"132": "53362031",
	"152": "1419215845",
	"162": "3272614282",
	"172": "668126235",
}

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

// GetRandomTrackFromDeezerGenre fetches a random track from a Deezer playlist
func GetRandomTrackFromDeezerGenre(genreID string) (*DeezerTrack, error) {
	playlistID, ok := PlaylistMappings[genreID]
	if !ok {
		playlistID = genreID
	}

	cacheMutex.RLock()
	tracks, found := playlistCache[playlistID]
	cacheMutex.RUnlock()

	if found && len(tracks) > 0 {
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(tracks))
		return &tracks[randomIndex], nil
	}

	var lastErr error
	for i := 0; i < 3; i++ {
		tracks, err := fetchTracksFromPlaylist(playlistID)
		if err == nil {
			cacheMutex.Lock()
			playlistCache[playlistID] = tracks
			cacheMutex.Unlock()

			if len(tracks) > 0 {
				rand.Seed(time.Now().UnixNano())
				randomIndex := rand.Intn(len(tracks))
				return &tracks[randomIndex], nil
			}
		}
		lastErr = err
		time.Sleep(500 * time.Millisecond)
	}

	return nil, fmt.Errorf("failed to fetch tracks after retries: %v", lastErr)
}

func fetchTracksFromPlaylist(playlistID string) ([]DeezerTrack, error) {
	url := fmt.Sprintf("%s/playlist/%s/tracks?limit=100", DeezerAPIURL, playlistID)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var allTracks []DeezerTrack

	// Fetch all pages of the playlist
	for url != "" {
		resp, err := client.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch playlist tracks: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %v", err)
		}

		var response DeezerPlaylistResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %v", err)
		}

		// Add tracks from this page
		for _, track := range response.Data {
			if track.Preview != "" {
				allTracks = append(allTracks, track)
			}
		}

		// Check if there are more pages (usually in a 'next' field for paginated responses)
		// For now, just get the first 100
		url = ""
	}

	if len(allTracks) == 0 {
		return nil, fmt.Errorf("playlist has no tracks with preview URLs")
	}

	return allTracks, nil
}
