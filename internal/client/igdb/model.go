package igdb

// TokenResp - access token response
type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// TopRatedGamesResp - top-rated games response
type TopRatedGamesResp struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	TotalRating      float64 `json:"total_rating"`
	TotalRatingCount int64   `json:"total_rating_count"`
	Cover            struct {
		URL string `json:"url"`
	} `json:"cover"`
	FirstReleaseDate int64 `json:"first_release_date"`
	Genres           []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	InvolvedCompanies []struct {
		Company struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		}
		Developer bool `json:"developer"`
		Publisher bool `json:"publisher"`
	} `json:"involved_companies"`
	Platforms   []int64 `json:"platforms"`
	Screenshots []struct {
		URL string `json:"url"`
	} `json:"screenshots"`
	Slug     string `json:"slug"`
	Summary  string `json:"summary"`
	Websites []struct {
		Category int8   `json:"category"`
		URL      string `json:"url"`
	} `json:"websites"`
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

// Image aliases
const (
	ImageThumbAlias         = "thumb"
	ImageLogoMed2xAlias     = "logo_med_2x"
	ImageScreenshotBigAlias = "screenshot_big"
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
