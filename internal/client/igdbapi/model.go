package igdbapi

import "bytes"

// TokenResp - access token response
type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// TopRatedGamesResp - top-rated games response
type TopRatedGamesResp struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	TotalRating       float64   `json:"total_rating"`
	TotalRatingCount  int64     `json:"total_rating_count"`
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
	Category int8   `json:"category"`
	URL      string `json:"url"`
}

// WebsiteCategory - website category
type WebsiteCategory int8

// Website categories
const (
	WebsiteCategoryOfficial  WebsiteCategory = 1
	WebsiteCategoryFacebook  WebsiteCategory = 4
	WebsiteCategoryTwitter   WebsiteCategory = 5
	WebsiteCategoryTwitch    WebsiteCategory = 6
	WebsiteCategoryYoutube   WebsiteCategory = 9
	WebsiteCategorySteam     WebsiteCategory = 13
	WebsiteCategoryEpicGames WebsiteCategory = 16
	WebsiteCategoryGOG       WebsiteCategory = 17
)

// Image types aliases
const (
	ImageTypeThumbAlias         = "thumb"
	ImageTypeCoverBig2xAlias    = "cover_big_2x"
	ImageTypeScreenshotBigAlias = "screenshot_big"
)

// WebsiteCategoryNames - mapping of website category to name
var WebsiteCategoryNames = map[WebsiteCategory]string{
	WebsiteCategoryOfficial:  "Official",
	WebsiteCategoryFacebook:  "Facebook",
	WebsiteCategoryTwitter:   "Twitter",
	WebsiteCategoryTwitch:    "Twitch",
	WebsiteCategoryYoutube:   "Youtube",
	WebsiteCategorySteam:     "Steam",
	WebsiteCategoryEpicGames: "EpicGames",
	WebsiteCategoryGOG:       "GOG",
}

// GetImageResp - get image response
type GetImageResp struct {
	Body        *bytes.Reader
	FileName    string
	ContentType string
}
