package model

import (
	"net/url"
	"strings"

	"cloud.google.com/go/civil"
	"golang.org/x/exp/constraints"
)

// Common validation errors
var (
	ErrRequiredMsg           = "required"
	ErrInvalidDateMsg        = "must be in format 'YYYY-MM-DD'"
	ErrInvalidRatingMsg      = "must be less or equal to 5"
	ErrInvalidImageURLMsg    = "must be a valid image URL from domain " + allowedImageDomains[0]
	ErrInvalidImageURLsMsg   = "must be valid images URLs from domain " + allowedImageDomains[0]
	ErrInvalidWebsitesURLMsg = "must be valid website URLs from domain list: " + strings.Join(allowedWebsiteDomains, ", ")
	ErrNonPositiveValuesMsg  = "must be greater than zero"
)

const (
	dateFieldLength = 10
)

// validates date format (YYYY-MM-DD)
func validateDate(date string) bool {
	if len(date) != dateFieldLength {
		return false
	}
	_, err := civil.ParseDate(date)
	return err == nil
}

// checks if URLs are from an allowed image CDN domain
func validateImageURLs(urls []string) bool {
	if len(urls) == 0 {
		return true
	}

	for _, imageURL := range urls {
		if !isValidURL(imageURL, allowedImageDomains) {
			return false
		}
	}
	return true
}

// checks if URLs are from allowed websites
func validateWebsiteURLs(urls []string) bool {
	if len(urls) == 0 {
		return true
	}

	for _, websiteURL := range urls {
		if !isValidURL(websiteURL, allowedWebsiteDomains) {
			return false
		}
	}
	return true
}

var allowedImageDomains = []string{"ucarecdn.com"}
var allowedWebsiteDomains = []string{
	"twitch.tv", "steampowered.com", "twitter.com", "facebook.com", "xbox.com",
	"youtube.com", "gog.com", "epicgames.com", "playstation.com"}

// checks if the URL is from an allowed domain list
func isValidURL(urlStr string, allowedList []string) bool {
	if urlStr == "" {
		return false
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil || parsedURL.Host == "" {
		return false
	}

	host := parsedURL.Host
	for _, domain := range allowedList {
		if strings.HasSuffix(host, domain) {
			return true
		}
	}

	return false
}

// checks if all slice values are positive
func validatePositive[T constraints.Integer](slice []T) bool {
	for _, v := range slice {
		if v <= 0 {
			return false
		}
	}
	return true
}

func removeDuplicates[T comparable](ids []T) []T {
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
