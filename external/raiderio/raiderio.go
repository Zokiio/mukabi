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
	"time"

	"github.com/topi314/tint"
)

// TODO: Refactor structs ito separate files for better organization
type RaiderIO struct {
	API      string `json:"api"`
	Key      string `json:"key"`
	Version  string `json:"version"`
	Cache    Cache  `json:"cache"`
	Timeout  int    `json:"timeout"`
	CacheMap *sync.Map
}

type Cache struct {
	Enabled bool   `json:"enabled"`
	TTL     int    `json:"ttl"`
	Backend string `json:"backend"`
}

type FilteredRealm struct {
	Region string `json:"region"`
	Realm  string `json:"realm"`
	Slug   string `json:"slug"`
}

type ConnectedRealms struct {
	RealmListing RealmListing `json:"realmListing"`
}

type RealmListing struct {
	Region    Region      `json:"region"`
	SubRegion interface{} `json:"subRegion"`
	Raid      Raid        `json:"raid"`
	Season    Season      `json:"season"`
	Realms    []Realm     `json:"realms"`
}
type Region struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	ShortName string `json:"short_name"`
}
type Raid struct {
	Type                     string      `json:"type"`
	ID                       int         `json:"id"`
	Difficulty               string      `json:"difficulty"`
	Name                     string      `json:"name"`
	ShortName                string      `json:"short_name"`
	IconURL                  string      `json:"icon_url"`
	Slug                     string      `json:"slug"`
	CanShowRaidMythicDetails bool        `json:"can_show_raid_mythic_details"`
	CanShowRaidHeroicDetails bool        `json:"can_show_raid_heroic_details"`
	CanShowRaidNormalDetails bool        `json:"can_show_raid_normal_details"`
	ExpansionID              int         `json:"expansion_id"`
	Encounters               []Encounter `json:"encounters"`
}
type Encounter struct {
	EncounterID int    `json:"encounterId"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Ordinal     int    `json:"ordinal"`
	WingID      int    `json:"wingId"`
	IconURL     string `json:"iconUrl"`
}
type Season struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}
type Realm struct {
	ID              int              `json:"id"`
	ConnectedRealms []ConnectedRealm `json:"connectedRealms"`
	Region          Region           `json:"region"`
	Stats           Stats            `json:"stats"`
}
type ConnectedRealm struct {
	Type     string      `json:"type"`
	Name     string      `json:"name"`
	AltName  interface{} `json:"alt_name"`
	Slug     string      `json:"slug"`
	Locale   string      `json:"locale"`
	Language string      `json:"language"`
	Timezone string      `json:"timezone"`
}
type Stats struct {
	NumAllianceCharacters      int   `json:"num_alliance_characters"`
	NumCombinedCharacters      int   `json:"num_combined_characters"`
	NumHordeCharacters         int   `json:"num_horde_characters"`
	NumAllianceGuilds          int   `json:"num_alliance_guilds"`
	NumCombinedGuilds          int   `json:"num_combined_guilds"`
	NumHordeGuilds             int   `json:"num_horde_guilds"`
	MplusHordeLevels           []int `json:"mplus_horde_levels"`
	MplusAllianceLevels        []int `json:"mplus_alliance_levels"`
	MplusCombinedLevels        []int `json:"mplus_combined_levels"`
	RaidCombinedNormalProgress []int `json:"raid_combined_normal_progress"`
	RaidCombinedHeroicProgress []int `json:"raid_combined_heroic_progress"`
	RaidCombinedMythicProgress []int `json:"raid_combined_mythic_progress"`
}

func New(ApiKey string) *RaiderIO {
	return &RaiderIO{
		API:     "https://raider.io/api",
		Key:     ApiKey,
		Version: "v1",
		Cache: Cache{
			Enabled: true,
			TTL:     3600, // 1 hour TTL
			Backend: "in-memory",
		},
		Timeout:  10,
		CacheMap: &sync.Map{},
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
func (r *RaiderIO) FetchConnectedRealms(region string, query string) ([]FilteredRealm, error) {
	// Cache key based on the region
	cacheKey := fmt.Sprintf("connected_realms_%s", region)

	// Check cache first
	cachedData, found := r.CacheMap.Load(cacheKey)
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
	url := fmt.Sprintf("%s/connected-realms?region=%s&realm=all", r.API, region)
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

	r.CacheMap.Store(cacheKey, dataToCache)

	// Apply query filtering before returning the results
	filteredRealms = filterRealms(filteredRealms, query)
	return filteredRealms, nil
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
	FieldMythicPlusScoresBySeason                 = "mythic_plus_scores_by_season"
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

type WoWCharacter struct {
	Name                 string               `json:"name"`
	Race                 string               `json:"race"`
	Class                string               `json:"class"`
	ActiveSpecName       string               `json:"active_spec_name"`
	ActiveSpecRole       string               `json:"active_spec_role"`
	Gender               string               `json:"gender"`
	Faction              string               `json:"faction"`
	AchievementPoints    int                  `json:"achievement_points"`
	ThumbnailURL         string               `json:"thumbnail_url"`
	Region               string               `json:"region"`
	Realm                string               `json:"realm"`
	LastCrawledAt        time.Time            `json:"last_crawled_at"`
	ProfileURL           string               `json:"profile_url"`
	ProfileBanner        string               `json:"profile_banner"`
	MythicPlusRecentRuns MythicPlusRecentRuns `json:"mythic_plus_recent_runs"`
}

type MythicPlusRecentRuns struct {
	Dungeon             string    `json:"dungeon"`
	ShortName           string    `json:"short_name"`
	MythicLevel         int       `json:"mythic_level"`
	CompletedAt         time.Time `json:"completed_at"`
	ClearTimeMs         int       `json:"clear_time_ms"`
	KeystoneRunID       int       `json:"keystone_run_id"`
	ParTimeMs           int       `json:"par_time_ms"`
	NumKeystoneUpgrades int       `json:"num_keystone_upgrades"`
	MapChallengeModeID  int       `json:"map_challenge_mode_id"`
	ZoneID              int       `json:"zone_id"`
	ZoneExpansionID     int       `json:"zone_expansion_id"`
	IconURL             string    `json:"icon_url"`
	BackgroundImageURL  string    `json:"background_image_url"`
	Score               int       `json:"score"`
	Affixes             Affixes   `json:"affixes"`
	URL                 string    `json:"url"`
}

type Affixes struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	IconURL     string `json:"icon_url"`
	WowheadURL  string `json:"wowhead_url"`
}

func (r *RaiderIO) FetchCharacterProfile(region, realm, character string, opts ...FetchCharacterOption) (*WoWCharacter, error) {
	endpoint := fmt.Sprintf("%s/%s/characters/profile", r.API, r.Version)
	params := fmt.Sprintf("?access_key=%s&region=%s&realm=%s&name=%s", r.Key, region, realm, character)

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

	var characterProfile WoWCharacter
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
