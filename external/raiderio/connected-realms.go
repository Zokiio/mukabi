// Package raiderio provides integration with the Raider.IO API
package raiderio

// ConnectedRealms represents the top-level response for connected realms data
type ConnectedRealms struct {
	RealmListing RealmListing `json:"realmListing"`
}

// RealmListing contains information about realms in a region
type RealmListing struct {
	Region    Region  `json:"region"`
	SubRegion any     `json:"subRegion"` // SubRegion data varies by region
	Raid      Raid    `json:"raid"`
	Season    Season  `json:"season"`
	Realms    []Realm `json:"realms"`
}

// Region represents a World of Warcraft game region
type Region struct {
	Name      string `json:"name"`       // Full region name
	Slug      string `json:"slug"`       // Region identifier
	ShortName string `json:"short_name"` // Abbreviated region name
}

// Raid represents current raid tier information
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

// Encounter represents a boss encounter within a raid
type Encounter struct {
	EncounterID int    `json:"encounterId"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Ordinal     int    `json:"ordinal"` // Boss order in raid
	WingID      int    `json:"wingId"`
	IconURL     string `json:"iconUrl"`
}

// Season represents a Mythic+ season
type Season struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// Realm represents a World of Warcraft realm and its connected realms
type Realm struct {
	ID              int              `json:"id"`
	ConnectedRealms []ConnectedRealm `json:"connectedRealms"`
	Region          Region           `json:"region"`
	Stats           Stats            `json:"stats"`
}

// ConnectedRealm represents a realm that shares a player pool with other realms
type ConnectedRealm struct {
	Type     string      `json:"type"`
	Name     string      `json:"name"`
	AltName  interface{} `json:"alt_name"` // Alternative realm name, if any
	Slug     string      `json:"slug"`     // Realm identifier
	Locale   string      `json:"locale"`   // Region locale
	Language string      `json:"language"` // Realm language
	Timezone string      `json:"timezone"` // Realm timezone
}

// Stats contains realm population and activity statistics
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
