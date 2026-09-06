package types_test

import (
	"testing"
	"time"

	"github.com/OutOfStack/game-library/pkg/types"
	"github.com/stretchr/testify/require"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  types.Date
		valid bool
	}{
		{name: "valid", input: "2026-09-05", want: types.Date{Year: 2026, Month: time.September, Day: 5}, valid: true},
		{name: "leap day", input: "2024-02-29", want: types.Date{Year: 2024, Month: time.February, Day: 29}, valid: true},
		{name: "empty"},
		{name: "non leap year", input: "2025-02-29"},
		{name: "invalid month", input: "2026-13-01"},
		{name: "invalid day", input: "2026-04-31"},
		{name: "wrong format", input: "05/09/2026"},
		{name: "timestamp", input: "2026-09-05T12:00:00Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := types.ParseDate(tt.input)
			if tt.valid {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.Error(t, err)
				require.Zero(t, got)
			}
		})
	}
}

func TestDate_String(t *testing.T) {
	tests := []struct {
		name string
		date types.Date
		want string
	}{
		{name: "zero"},
		{name: "padded components", date: types.Date{Year: 7, Month: time.March, Day: 2}, want: "0007-03-02"},
		{name: "leap day", date: types.Date{Year: 2024, Month: time.February, Day: 29}, want: "2024-02-29"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.date.String())
		})
	}
}

func TestDate_Scan(t *testing.T) {
	initial := types.Date{Year: 2000, Month: time.January, Day: 1}
	tests := []struct {
		name  string
		input any
		want  types.Date
		valid bool
	}{
		{name: "nil clears date", valid: true},
		{
			name: "local calendar date", valid: true,
			input: time.Date(2024, time.March, 1, 0, 30, 0, 0, time.FixedZone("UTC+3", 3*60*60)),
			want:  types.Date{Year: 2024, Month: time.March, Day: 1},
		},
		{name: "string rejected", input: "2026-09-05", want: initial},
		{name: "integer rejected", input: 42, want: initial},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := initial
			err := got.Scan(tt.input)
			if tt.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDateOf(t *testing.T) {
	input := time.Date(2024, time.March, 1, 0, 30, 59, 123, time.FixedZone("UTC+3", 3*60*60))
	require.Equal(t, types.Date{Year: 2024, Month: time.March, Day: 1}, types.DateOf(input))
	require.Equal(t, types.Date{Year: 2024, Month: time.February, Day: 29}, types.DateOf(input.UTC()))
}
