package model_test

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/model"
	"github.com/stretchr/testify/require"
)

func TestTaskSettings_Value(t *testing.T) {
	for _, input := range []model.TaskSettings{nil, {}, []byte(`{"enabled":true}`)} {
		t.Run(string(input), func(t *testing.T) {
			value, err := input.Value()
			require.NoError(t, err)
			require.Equal(t, string(input), value)
			var got model.TaskSettings
			require.NoError(t, got.Scan(value))
			require.Equal(t, string(input), string(got))
		})
	}
}

func TestTaskSettings_Scan(t *testing.T) {
	const initial = `{"original":true}`
	tests := []struct {
		name    string
		input   any
		want    model.TaskSettings
		wantErr bool
	}{
		{name: "string", input: `{"enabled":true}`, want: model.TaskSettings(`{"enabled":true}`)},
		{name: "bytes", input: []byte(`{"count":2}`), want: model.TaskSettings(`{"count":2}`)},
		{name: "nil preserves value", want: model.TaskSettings(initial)},
		{name: "nil bytes", input: []byte(nil)},
		{name: "empty string", input: "", want: model.TaskSettings{}},
		{name: "unsupported type", input: 42, want: model.TaskSettings(initial), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := model.TaskSettings(initial)
			err := got.Scan(tt.input)
			if tt.wantErr {
				require.ErrorContains(t, err, "unsupported type int")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
