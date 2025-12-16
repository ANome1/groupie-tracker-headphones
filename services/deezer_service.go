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

// Cache for genres to avoid hitting API every round
var (
	genreCache = make(map[string][]DeezerTrack)
	cacheMutex sync.RWMutex
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

// GetRandomTrackFromDeezerGenre fetches a random track from a Deezer genre
func GetRandomTrackFromDeezerGenre(genreID string) (*DeezerTrack, error) {
	// Check cache first
	cacheMutex.RLock()
	tracks, found := genreCache[genreID]
	cacheMutex.RUnlock()

	if found && len(tracks) > 0 {
		rand.Seed(time.Now().UnixNano())
		randomIndex := rand.Intn(len(tracks))
		return &tracks[randomIndex], nil
	}

	// Not in cache, fetch from API
	var lastErr error
	for i := 0; i < 3; i++ {
		tracks, err := fetchTracksFromGenre(genreID)
		if err == nil {
			// Update cache
			cacheMutex.Lock()
			genreCache[genreID] = tracks
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

func fetchTracksFromGenre(genreID string) ([]DeezerTrack, error) {
	// Use Deezer radio endpoint for genre (always works)
	url := fmt.Sprintf("%s/radio/%s/tracks?limit=100", DeezerAPIURL, genreID)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch genre tracks: %v", err)
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

	// Filter tracks to keep only those with valid preview URLs
	var validTracks []DeezerTrack
	for _, track := range response.Data {
		if track.Preview != "" {
			validTracks = append(validTracks, track)
		}
	}

	if len(validTracks) == 0 {
		return nil, fmt.Errorf("genre has no tracks with preview URLs")
	}

	return validTracks, nil
}

// ClearCache clears the genre cache
func ClearDeezerCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	genreCache = make(map[string][]DeezerTrack)
}
