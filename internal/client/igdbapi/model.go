package igdbapi

import "bytes"

// Image types aliases
const (
	ImageTypeThumbAlias         = "thumb"
	ImageTypeCoverBig2xAlias    = "cover_big_2x"
	ImageTypeScreenshotBigAlias = "screenshot_big"
)

// TokenResp - access token response
type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// TopRatedGames - top-rated games
type TopRatedGames struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	TotalRating       float64   `json:"total_rating"`
	TotalRatingCount  int32     `json:"total_rating_count"`
	Cover             URL       `json:"cover"`
	FirstReleaseDate  int64     `json:"first_release_date"`
	Genres            []IDName  `json:"genres"`
	InvolvedCompanies []Company `json:"involved_companies"`
	Platforms         []int64   `json:"platforms"`
	Screenshots       []URL     `json:"screenshots"`
	Slug              string    `json:"slug"`
	Summary           string    `json:"summary"`
	Websites          []Website `json:"websites"`
}

// URL - struct containing url
type URL struct {
	URL string `json:"url"`
}

// IDName - struct containing id and name
type IDName struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Company - company
type Company struct {
	Company   IDName `json:"company"`
	Developer bool   `json:"developer"`
	Publisher bool   `json:"publisher"`
}

// Website - website
type Website struct {
	URL  string `json:"url"`
	Type int8   `json:"type"`
}

// Website types
const (
	WebsiteTypeOfficial  int8 = 1
	WebsiteTypeFacebook  int8 = 4
	WebsiteTypeTwitter   int8 = 5
	WebsiteTypeTwitch    int8 = 6
	WebsiteTypeYoutube   int8 = 9
	WebsiteTypeSteam     int8 = 13
	WebsiteTypeEpicGames int8 = 16
	WebsiteTypeGOG       int8 = 17
	WebsiteTypeDiscord   int8 = 18
)

// WebsiteTypeNames - mapping of a website type to name
var WebsiteTypeNames = map[int8]string{
	WebsiteTypeOfficial:  "Official",
	WebsiteTypeFacebook:  "Facebook",
	WebsiteTypeTwitter:   "Twitter",
	WebsiteTypeTwitch:    "Twitch",
	WebsiteTypeYoutube:   "Youtube",
	WebsiteTypeSteam:     "Steam",
	WebsiteTypeEpicGames: "EpicGames",
	WebsiteTypeGOG:       "GOG",
	WebsiteTypeDiscord:   "Discord",
}

// GetImageResp - get image response
type GetImageResp struct {
	Body        *bytes.Reader
	FileName    string
	ContentType string
}

// GameInfoForUpdate - game info for update
type GameInfoForUpdate struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	TotalRating      float64   `json:"total_rating"`
	TotalRatingCount int32     `json:"total_rating_count"`
	Platforms        []int64   `json:"platforms"`
	Websites         []Website `json:"websites"`
}
