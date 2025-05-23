package validation

import (
	"net/url"
	"strings"

	"cloud.google.com/go/civil"
	"github.com/OutOfStack/game-library/internal/appconf"
	"go.uber.org/zap"
)

const (
	dateFieldLength = 10
)

var allowedWebsiteDomains = []string{
	"twitch.tv", "steampowered.com", "twitter.com", "facebook.com", "xbox.com",
	"youtube.com", "gog.com", "epicgames.com", "playstation.com"}

// Validator struct
type Validator struct {
	allowedImageDomain string
}

// NewValidator creates new validator
func NewValidator(log *zap.Logger, cfg *appconf.Cfg) *Validator {
	var allowedImageDomain string
	u, err := url.Parse(cfg.GetS3().CDNBaseURL)
	if err != nil {
		// url is validated on config file/env read
		// if we still failed to provide correct domain, allow any
		log.Error("can't parse CDN base URL, allow any in validation", zap.String("url", cfg.GetS3().CDNBaseURL), zap.Error(err))
	} else {
		allowedImageDomain = u.Host
	}

	return &Validator{allowedImageDomain: allowedImageDomain}
}

// Common validation errors

// ErrRequiredMsg returns error message
func (v *Validator) ErrRequiredMsg() string {
	return "required"
}

// ErrInvalidDateMsg returns error message
func (v *Validator) ErrInvalidDateMsg() string {
	return "must be in format 'YYYY-MM-DD'"
}

// ErrInvalidRatingMsg returns error message
func (v *Validator) ErrInvalidRatingMsg() string {
	return "must be less or equal to 5"
}

// ErrInvalidImageURLMsg returns error message
func (v *Validator) ErrInvalidImageURLMsg() string {
	return "must be a valid image URL from domain " + v.allowedImageDomain
}

// ErrInvalidImageURLsMsg returns error message
func (v *Validator) ErrInvalidImageURLsMsg() string {
	return "must be valid images URLs from domain " + v.allowedImageDomain
}

// ErrInvalidWebsitesURLMsg returns error message
func (v *Validator) ErrInvalidWebsitesURLMsg() string {
	return "must be valid website URLs from domain list: " + strings.Join(allowedWebsiteDomains, ", ")
}

// ErrNonPositiveValuesMsg returns error message
func (v *Validator) ErrNonPositiveValuesMsg() string {
	return "must be greater than zero"
}

// ValidateDate validates date format (YYYY-MM-DD)
func (v *Validator) ValidateDate(date string) bool {
	if len(date) != dateFieldLength {
		return false
	}
	_, err := civil.ParseDate(date)
	return err == nil
}

// ValidateImageURLs checks if URLs are from an allowed image CDN domain
func (v *Validator) ValidateImageURLs(urls []string) bool {
	if len(urls) == 0 {
		return true
	}

	for _, imageURL := range urls {
		if !v.IsValidURL(imageURL, []string{v.allowedImageDomain}) {
			return false
		}
	}
	return true
}

// ValidateWebsiteURLs checks if URLs are from allowed websites
func (v *Validator) ValidateWebsiteURLs(urls []string) bool {
	if len(urls) == 0 {
		return true
	}

	for _, websiteURL := range urls {
		if !v.IsValidURL(websiteURL, allowedWebsiteDomains) {
			return false
		}
	}
	return true
}

// IsValidURL checks if the URL is from an allowed domain list
func (v *Validator) IsValidURL(urlStr string, allowedList []string) bool {
	if urlStr == "" {
		return false
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil || parsedURL.Host == "" {
		return false
	}

	host := parsedURL.Host
	for _, allowed := range allowedList {
		if strings.HasSuffix(host, allowed) {
			return true
		}
	}

	return false
}

// ValidatePositive checks if all slice values are positive
func (v *Validator) ValidatePositive(slice []int32) bool {
	for _, v := range slice {
		if v <= 0 {
			return false
		}
	}
	return true
}

// RemoveDuplicates removes duplicates from slice
func RemoveDuplicates[T comparable](ids []T) []T {
	unique := make(map[T]bool)
	var res []T
	for _, id := range ids {
		if !unique[id] {
			unique[id] = true
			res = append(res, id)
		}
	}
	return res
}
