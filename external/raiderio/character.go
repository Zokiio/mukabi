package raiderio

import "time"

// CharacterProfile represents the main profile information.
type CharacterProfile struct {
	Name              string         `json:"name"`
	Race              string         `json:"race"`
	Class             string         `json:"class"`
	ActiveSpecName    string         `json:"active_spec_name"`
	ActiveSpecRole    string         `json:"active_spec_role"`
	Gender            string         `json:"gender"`
	Faction           string         `json:"faction"`
	AchievementPoints int            `json:"achievement_points"`
	ThumbnailURL      string         `json:"thumbnail_url"`
	Region            string         `json:"region"`
	Realm             string         `json:"realm"`
	LastCrawledAt     time.Time      `json:"last_crawled_at"`
	ProfileURL        string         `json:"profile_url"`
	ProfileBanner     string         `json:"profile_banner"`
	Covenant          any            `json:"covenant"` // Consider creating a specific Covenant struct if the any holds consistent data
	Guild             *Guild         `json:"guild"`
	TalentData        *TalentLoadout `json:"talentLoadout"`
}

// Guild represents the guild information.
type Guild struct {
	Name  string `json:"name"`
	Realm string `json:"realm"`
}

// MythicPlusData encompasses all Mythic+ related information.
type MythicPlusData struct {
	ScoresBySeason []any                   `json:"mythic_plus_scores_by_season"`
	Ranks          MythicPlusRanks         `json:"mythic_plus_ranks"`
	PreviousRanks  PreviousMythicPlusRanks `json:"previous_mythic_plus_ranks"`
	RecentRuns     []MythicPlusRun         `json:"mythic_plus_recent_runs"`
	BestRuns       []MythicPlusRun         `json:"mythic_plus_best_runs"`
	HighestRuns    []MythicPlusRun         `json:"mythic_plus_highest_level_runs"`
	WeeklyHighest  []MythicPlusRun         `json:"mythic_plus_weekly_highest_level_runs"`
	PreviousWeekly []MythicPlusRun         `json:"mythic_plus_previous_weekly_highest_level_runs"`
}

// MythicPlusRanks represents the various Mythic+ ranks.
type MythicPlusRanks struct {
	Overall     Rank `json:"overall"`
	Class       Rank `json:"class"`
	Healer      Rank `json:"healer"`
	ClassHealer Rank `json:"class_healer"`
	Dps         Rank `json:"dps"`
	ClassDps    Rank `json:"class_dps"`
	Spec256     Rank `json:"spec_256"`
	Spec257     Rank `json:"spec_257"`
	Spec258     Rank `json:"spec_258"`
}

// PreviousMythicPlusRanks represents the previous Mythic+ ranks.
type PreviousMythicPlusRanks struct {
	Overall     Rank `json:"overall"`
	Class       Rank `json:"class"`
	Healer      Rank `json:"healer"`
	ClassHealer Rank `json:"class_healer"`
	Dps         Rank `json:"dps"`
	ClassDps    Rank `json:"class_dps"`
	Spec256     Rank `json:"spec_256"`
	Spec257     Rank `json:"spec_257"`
	Spec258     Rank `json:"spec_258"`
}

// Rank represents the world, region, and realm rank.
type Rank struct {
	World  int `json:"world"`
	Region int `json:"region"`
	Realm  int `json:"realm"`
}

// MythicPlusRun represents a single Mythic+ dungeon run.
type MythicPlusRun struct {
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
	Score               float64   `json:"score"`
	Affixes             []Affix   `json:"affixes"`
	URL                 string    `json:"url"`
}

// Affix represents a Mythic+ affix.
type Affix struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	IconURL     string `json:"icon_url"`
	WowheadURL  string `json:"wowhead_url"`
}

// Gear represents the character's equipped gear.
type Gear struct {
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	Source            string     `json:"source"`
	ItemLevelEquipped int        `json:"item_level_equipped"`
	ItemLevelTotal    int        `json:"item_level_total"`
	ArtifactTraits    int        `json:"artifact_traits"`
	Corruption        Corruption `json:"corruption"`
	Items             Items      `json:"items"`
}

// Corruption represents the character's corruption stats.
type Corruption struct {
	Added     int   `json:"added"`
	Resisted  int   `json:"resisted"`
	Total     int   `json:"total"`
	CloakRank int   `json:"cloakRank"`
	Spells    []any `json:"spells"`
}

// Items represents the equipped items in each slot.
type Items struct {
	Head     *Item `json:"head"`
	Neck     *Item `json:"neck"`
	Shoulder *Item `json:"shoulder"`
	Back     *Item `json:"back"`
	Chest    *Item `json:"chest"`
	Waist    *Item `json:"waist"`
	Wrist    *Item `json:"wrist"`
	Hands    *Item `json:"hands"`
	Legs     *Item `json:"legs"`
	Feet     *Item `json:"feet"`
	Finger1  *Item `json:"finger1"`
	Finger2  *Item `json:"finger2"`
	Trinket1 *Item `json:"trinket1"`
	Trinket2 *Item `json:"trinket2"`
	Mainhand *Item `json:"mainhand"`
	Offhand  *Item `json:"offhand"`
}

// Item represents a single equipped item.
type Item struct {
	ItemID           int             `json:"item_id"`
	ItemLevel        int             `json:"item_level"`
	Icon             string          `json:"icon"`
	Name             string          `json:"name"`
	ItemQuality      int             `json:"item_quality"`
	IsLegendary      bool            `json:"is_legendary"`
	IsAzeriteArmor   bool            `json:"is_azerite_armor"`
	AzeritePowers    []any           `json:"azerite_powers"`
	Corruption       *ItemCorruption `json:"corruption"`
	DominationShards []any           `json:"domination_shards"`
	Gems             []any           `json:"gems"`
	Enchants         []any           `json:"enchants"`
	Bonuses          []int           `json:"bonuses"`
	Enchant          int             `json:"enchant,omitempty"`
}

// ItemCorruption represents the corruption on a specific item.
type ItemCorruption struct {
	Added    int `json:"added"`
	Resisted int `json:"resisted"`
	Total    int `json:"total"`
}

// TalentLoadout represents the character's talent loadout.
type TalentLoadout struct {
	LoadoutSpecID  int               `json:"loadout_spec_id"`
	LoadoutText    string            `json:"loadout_text"`
	Loadout        []TalentNodeEntry `json:"loadout"`
	ClassTalents   []TalentNodeEntry `json:"class_talents"`
	SpecTalents    []TalentNodeEntry `json:"spec_talents"`
	HeroTalents    []TalentNodeEntry `json:"hero_talents"`
	ActiveHeroTree *ActiveHeroTree   `json:"active_hero_tree"`
}

// TalentNodeEntry represents a single talent node and its rank.
type TalentNodeEntry struct {
	Node             TalentNode `json:"node"`
	EntryIndex       int        `json:"entryIndex"`
	Rank             int        `json:"rank"`
	IncludeInSummary bool       `json:"includeInSummary,omitempty"`
}

// TalentNode represents a node in the talent tree.
type TalentNode struct {
	ID        int           `json:"id"`
	TreeID    int           `json:"treeId"`
	SubTreeID int           `json:"subTreeId"`
	Type      int           `json:"type"`
	Entries   []TalentEntry `json:"entries"`
	Important bool          `json:"important"`
	PosX      int           `json:"posX"`
	PosY      int           `json:"posY"`
	Row       int           `json:"row"`
	Col       int           `json:"col"`
}

// TalentEntry represents a specific talent within a node.
type TalentEntry struct {
	ID                int         `json:"id"`
	TraitDefinitionID int         `json:"traitDefinitionId"`
	TraitSubTreeID    int         `json:"traitSubTreeId"`
	Type              int         `json:"type"`
	MaxRanks          int         `json:"maxRanks"`
	Spell             TalentSpell `json:"spell"`
}

// TalentSpell represents the spell associated with a talent.
type TalentSpell struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	School      int    `json:"school"`
	Rank        any    `json:"rank"`
	HasCooldown bool   `json:"hasCooldown"`
}

// ActiveHeroTree represents the currently active hero talent tree.
type ActiveHeroTree struct {
	ID          int    `json:"id"`
	TraitTreeID int    `json:"traitTreeId"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	IconURL     string `json:"iconUrl"`
}

// RaidProgression represents the character's raid progression.
type RaidProgression struct {
	LiberationOfUndermine RaidEncounter `json:"liberation-of-undermine"`
	BlackrockDepths       RaidEncounter `json:"blackrock-depths"`
	NerubarPalace         RaidEncounter `json:"nerubar-palace"`
}

// RaidEncounter represents the progression in a specific raid encounter.
type RaidEncounter struct {
	Summary            string `json:"summary"`
	ExpansionID        int    `json:"expansion_id"`
	TotalBosses        int    `json:"total_bosses"`
	NormalBossesKilled int    `json:"normal_bosses_killed"`
	HeroicBossesKilled int    `json:"heroic_bosses_killed"`
	MythicBossesKilled int    `json:"mythic_bosses_killed"`
}

// CharacterAchievements encompasses raid achievement data.
type CharacterAchievements struct {
	Meta  []any `json:"raid_achievement_meta"`  // Consider creating a specific AchievementMeta struct if the any holds consistent data
	Curve []any `json:"raid_achievement_curve"` // Consider creating a specific AchievementCurve struct if the any holds consistent data
}

// CharacterProfile is the main struct that holds all the organized data.
type CharacterProfileResponse struct {
	Character    CharacterProfile      `json:"character_profile"`
	MythicPlus   MythicPlusData        `json:"mythic_plus_data"`
	Gear         Gear                  `json:"gear"`
	Talents      TalentLoadout         `json:"talent_loadout"`
	RaidProg     RaidProgression       `json:"raid_progression"`
	Achievements CharacterAchievements `json:"achievements"`
}
