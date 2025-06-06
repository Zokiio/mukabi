// Package raiderio provides integration with the Raider.IO API for World of Warcraft character and realm data.
package raiderio

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/topi314/tint"
)

const (
	defaultAPIURL     = "https://raider.io/api"
	defaultAPIVersion = "v1"
	defaultCacheTTL   = 3600 // 1 hour TTL
)

// Client represents a RaiderIO API client with caching capabilities.
type Client struct {
	apiURL     string
	apiKey     string
	apiVersion string
	cache      cacheConfig
	timeout    int
	cacheStore *sync.Map
}

type cacheConfig struct {
	enabled bool
	ttl     int
	backend string
}

// New creates a new RaiderIO client with the given API key.
func New(apiKey string) *Client {
	return &Client{
		apiURL:     defaultAPIURL,
		apiKey:     apiKey,
		apiVersion: defaultAPIVersion,
		cache: cacheConfig{
			enabled: true,
			ttl:     defaultCacheTTL,
			backend: "in-memory",
		},
		timeout:    10,
		cacheStore: &sync.Map{},
	}
}

// filterRealms filters the realms based on the provided query string
// It returns a slice of FilteredRealm that match the query
func filterRealms(realms []FilteredRealm, query string) []FilteredRealm {
	var filtered []FilteredRealm

	// If query is empty, return all realms
	if query == "" {
		return realms
	}

	// Apply the filter: match realms where the realm name contains the query
	for _, realm := range realms {
		if strings.Contains(strings.ToLower(realm.Realm), strings.ToLower(query)) {
			filtered = append(filtered, realm)
		}
	}

	return filtered
}

// FetchConnectedRealms fetches connected realms from the RaiderIO API and caches the result
// It filters the realms based on the provided query string and returns the filtered list
func (c *Client) FetchConnectedRealms(region string, query string) ([]FilteredRealm, error) {
	// Cache key based on the region
	cacheKey := fmt.Sprintf("connected_realms_%s", region)

	// Check cache first
	cachedData, found := c.cacheStore.Load(cacheKey)
	if found {
		// If cache hit, unmarshal and return the filtered realms
		var realms []FilteredRealm
		if err := json.Unmarshal(cachedData.([]byte), &realms); err != nil {
			slog.Error("Failed to unmarshal cached data", tint.Err(err))
			return nil, fmt.Errorf("error unmarshaling cached data: %w", err)
		}
		// Apply query filtering
		filteredRealms := filterRealms(realms, query)
		return filteredRealms, nil
	}

	// If no cache hit, fetch data from the API
	url := fmt.Sprintf("%s/connected-realms?region=%s&realm=all", c.apiURL, region)
	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Failed to make API request", tint.Err(err))
		return nil, fmt.Errorf("error making API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Received non-OK status code", slog.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("error: received non-OK status code %d", resp.StatusCode)
	}

	// Decode API response
	var apiResponse ConnectedRealms
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		slog.Error("Failed to decode API response", tint.Err(err))
		return nil, fmt.Errorf("error decoding API response: %w", err)
	}

	// Map the realms to FilteredRealm
	var filteredRealms []FilteredRealm
	for _, realm := range apiResponse.RealmListing.Realms {
		for _, connectedRealm := range realm.ConnectedRealms {
			filteredRealms = append(filteredRealms, FilteredRealm{
				Region: region,
				Realm:  connectedRealm.Name,
				Slug:   connectedRealm.Slug,
			})
		}
	}

	// Cache the fetched data
	dataToCache, err := json.Marshal(filteredRealms)
	if err != nil {
		slog.Error("Failed to marshal data for cache", tint.Err(err))
		return nil, fmt.Errorf("error marshaling data for cache: %w", err)
	}

	c.cacheStore.Store(cacheKey, dataToCache)

	filteredRealms = filterRealms(filteredRealms, query)
	return filteredRealms, nil
}

type FilteredRealm struct {
	Region string `json:"region"`
	Realm  string `json:"realm"`
	Slug   string `json:"slug"`
}

type CharacterProfileReq struct {
	AccessKey string `json:"access_key"`
	Region    string `json:"region"`
	Realm     string `json:"realm"`
	Character string `json:"character"`
}

const (
	FieldGear                                     = "gear"
	FieldTalents                                  = "talents"
	FieldTalentsCategorized                       = "talents:categorized"
	FieldGuild                                    = "guild"
	FieldCovenant                                 = "covenant"
	FieldRaidProgression                          = "raid_progression"
	FieldMythicPlusScoresBySeason                 = "mythic_plus_scores_by_season:current"
	FieldMythicPlusRanks                          = "mythic_plus_ranks"
	FieldMythicPlusRecentRuns                     = "mythic_plus_recent_runs"
	FieldMythicPlusBestRuns                       = "mythic_plus_best_runs"
	FieldMythicPlusAlternateRuns                  = "mythic_plus_alternate_runs"
	FieldMythicPlusHighestLevelRuns               = "mythic_plus_highest_level_runs"
	FieldMythicPlusWeeklyHighestLevelRuns         = "mythic_plus_weekly_highest_level_runs"
	FieldMythicPlusPreviousWeeklyHighestLevelRuns = "mythic_plus_previous_weekly_highest_level_runs"
	FieldPreviousMythicPlusRanks                  = "previous_mythic_plus_ranks"
	FieldRaidAchievementMeta                      = "raid_achievement_meta"
	FieldRaidAchievementCurve                     = "raid_achievement_curve"
)

type FetchCharacterOption func(*fetchCharacterConfig)

type fetchCharacterConfig struct {
	fields []string
}

func WithFields(fields ...string) FetchCharacterOption {
	return func(cfg *fetchCharacterConfig) {
		cfg.fields = append(cfg.fields, fields...)
	}
}

func (c *Client) FetchCharacterProfile(region, realm, character string, opts ...FetchCharacterOption) (*CharacterProfile, error) {
	endpoint := fmt.Sprintf("%s/%s/characters/profile", c.apiURL, c.apiVersion)
	params := fmt.Sprintf("?access_key=%s&region=%s&realm=%s&name=%s", c.apiKey, region, realm, character)

	cfg := &fetchCharacterConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if len(cfg.fields) > 0 {
		params += "&fields=" + strings.Join(cfg.fields, ",")
	}

	fullURL := endpoint + params
	slog.Debug("Fetching character profile from URL", slog.String("url", fullURL))

	resp, err := http.Get(fullURL)
	if err != nil {
		slog.Error("Failed to make API request", tint.Err(err))
		return nil, fmt.Errorf("error making API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Received non-OK status code", slog.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", tint.Err(err))
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var characterProfile CharacterProfile
	if err := json.Unmarshal(body, &characterProfile); err != nil {
		slog.Error("Failed to unmarshal character profile", tint.Err(err))
		return nil, fmt.Errorf("error unmarshaling character profile: %w", err)
	}

	characterProfile.Gender = strings.ToTitle(characterProfile.Gender)
	if characterProfile.ThumbnailURL != "" {
		if _, err := url.ParseRequestURI(characterProfile.ThumbnailURL); err != nil {
			slog.Warn("Invalid thumbnail URL", slog.String("url", characterProfile.ThumbnailURL))
			characterProfile.ThumbnailURL = "" // Clear it to prevent invalid embeds
		}
	}
	if characterProfile.ProfileBanner != "" {
		if _, err := url.ParseRequestURI(characterProfile.ProfileBanner); err != nil {
			slog.Warn("Invalid thumbnail URL", slog.String("url", characterProfile.ProfileBanner))
			characterProfile.ProfileBanner = "" // Clear it to prevent invalid embeds
		}
	}

	return &characterProfile, nil
}
