package validation

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/stretchr/testify/assert"
)

func TestValidateDate(t *testing.T) {
	v := NewValidator(nil, &appconf.Cfg{})
	tests := []struct {
		name     string
		date     string
		expected bool
	}{
		{"valid date", "2025-05-22", true},
		{"invalid format", "2025-05-2", false},
		{"invalid year", "2025-13-01", false},
		{"too short", "2025-05", false},
		{"too long", "2025-05-22-01", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, v.ValidateDate(tt.date))
		})
	}
}

func TestValidateImageURLs(t *testing.T) {
	cfg := &appconf.Cfg{
		S3: appconf.S3{
			CDNBaseURL: "https://cdn.example.com",
		},
	}
	v := NewValidator(nil, cfg)

	tests := []struct {
		name     string
		urls     []string
		expected bool
	}{
		{"empty slice", []string{}, true},
		{"valid URLs", []string{"https://cdn.example.com/image1.jpg", "https://cdn.example.com/image2.jpg"}, true},
		{"invalid domain", []string{"https://other-domain.com/image.jpg"}, false},
		{"mixed URLs", []string{"https://cdn.example.com/image.jpg", "https://other-domain.com/image.jpg"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, v.ValidateImageURLs(tt.urls))
		})
	}
}

func TestValidateWebsiteURLs(t *testing.T) {
	v := NewValidator(nil, &appconf.Cfg{})
	tests := []struct {
		name     string
		urls     []string
		expected bool
	}{
		{"empty slice", []string{}, true},
		{"valid URLs", []string{"https://twitch.tv/user", "https://twitter.com/user"}, true},
		{"invalid domain", []string{"https://invalid-domain.com"}, false},
		{"mixed URLs", []string{"https://twitch.tv/user", "https://invalid-domain.com"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, v.ValidateWebsiteURLs(tt.urls))
		})
	}
}

func TestValidatePositive(t *testing.T) {
	v := NewValidator(nil, &appconf.Cfg{})
	tests := []struct {
		name     string
		slice    []int32
		expected bool
	}{
		{"all positive", []int32{1, 2, 3}, true},
		{"contains zero", []int32{1, 0, 3}, false},
		{"contains negative", []int32{1, -1, 3}, false},
		{"empty slice", []int32{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, v.ValidatePositive(tt.slice))
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		input    []int32
		expected []int32
	}{
		{"no duplicates", []int32{1, 2, 3}, []int32{1, 2, 3}},
		{"with duplicates", []int32{1, 2, 2, 3, 3, 3}, []int32{1, 2, 3}},
		{"empty slice", []int32{}, nil},
		{"all duplicates", []int32{1, 1, 1}, []int32{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, RemoveDuplicates(tt.input))
		})
	}
}
