// Package raiderio provides integration with the Raider.IO API
package raiderio

import "time"

// CharacterProfile represents detailed information about a World of Warcraft character
type CharacterProfile struct {
	Name                     string                     `json:"name"`
	Race                     string                     `json:"race"`
	Class                    string                     `json:"class"`
	ActiveSpecName           string                     `json:"active_spec_name"`
	ActiveSpecRole           string                     `json:"active_spec_role"`
	Gender                   string                     `json:"gender"`
	Faction                  string                     `json:"faction"`
	AchievementPoints        int                        `json:"achievement_points"`
	ThumbnailURL             string                     `json:"thumbnail_url"`
	Region                   string                     `json:"region"`
	Realm                    string                     `json:"realm"`
	LastCrawledAt            time.Time                  `json:"last_crawled_at"`
	ProfileURL               string                     `json:"profile_url"`
	ProfileBanner            string                     `json:"profile_banner"`
	MythicPlusScoresBySeason []MythicPlusScoresBySeason `json:"mythic_plus_scores_by_season"`
}

// MythicPlusScoresBySeason represents Mythic+ scores for a specific season
type MythicPlusScoresBySeason struct {
	Season   string             `json:"season"`   // Season identifier
	Scores   MythicPlusScores   `json:"scores"`   // Overall and role-specific scores
	Segments MythicPlusSegments `json:"segments"` // Detailed score breakdowns
}

// MythicPlusScores contains various Mythic+ rating scores for different roles
type MythicPlusScores struct {
	All    float64 `json:"all"`    // Overall Mythic+ score
	Dps    int     `json:"dps"`    // DPS role score
	Healer float64 `json:"healer"` // Healer role score
	Tank   int     `json:"tank"`   // Tank role score
	Spec0  float64 `json:"spec_0"` // First specialization score
	Spec1  int     `json:"spec_1"` // Second specialization score
	Spec2  int     `json:"spec_2"` // Third specialization score
	Spec3  int     `json:"spec_3"` // Fourth specialization score
}

// MythicPlusSegments contains detailed score information for different roles
type MythicPlusSegments struct {
	All    MythicPlusSegment `json:"all"`    // Overall segment data
	Dps    MythicPlusSegment `json:"dps"`    // DPS role segment
	Healer MythicPlusSegment `json:"healer"` // Healer role segment
	Tank   MythicPlusSegment `json:"tank"`   // Tank role segment
	Spec0  MythicPlusSegment `json:"spec_0"` // First spec segment
	Spec1  MythicPlusSegment `json:"spec_1"` // Second spec segment
	Spec2  MythicPlusSegment `json:"spec_2"` // Third spec segment
	Spec3  MythicPlusSegment `json:"spec_3"` // Fourth spec segment
}

// MythicPlusSegment represents a score segment with color coding
type MythicPlusSegment struct {
	Score float64 `json:"score"` // Segment score value
	Color string  `json:"color"` // Color code for the score range
}
