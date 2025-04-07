package raiderio

import "time"

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

type MythicPlusScoresBySeason struct {
	Season   string             `json:"season"`
	Scores   MythicPlusScores   `json:"scores"`
	Segments MythicPlusSegments `json:"segments"`
}

type MythicPlusScores struct {
	All    float64 `json:"all"`
	Dps    int     `json:"dps"`
	Healer float64 `json:"healer"`
	Tank   int     `json:"tank"`
	Spec0  float64 `json:"spec_0"`
	Spec1  int     `json:"spec_1"`
	Spec2  int     `json:"spec_2"`
	Spec3  int     `json:"spec_3"`
}

type MythicPlusSegments struct {
	All    MythicPlusSegment `json:"all"`
	Dps    MythicPlusSegment `json:"dps"`
	Healer MythicPlusSegment `json:"healer"`
	Tank   MythicPlusSegment `json:"tank"`
	Spec0  MythicPlusSegment `json:"spec_0"`
	Spec1  MythicPlusSegment `json:"spec_1"`
	Spec2  MythicPlusSegment `json:"spec_2"`
	Spec3  MythicPlusSegment `json:"spec_3"`
}

type MythicPlusSegment struct {
	Score float64 `json:"score"`
	Color string  `json:"color"`
}
