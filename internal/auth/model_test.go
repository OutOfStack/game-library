package auth_test

import (
	"testing"

	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
)

func TestClaims_UserID(t *testing.T) {
	subject := td.String()
	claims := &auth.Claims{Name: td.String()}
	claims.Subject = subject
	tests := []struct {
		name   string
		claims *auth.Claims
		want   string
	}{
		{name: "nil claims"},
		{name: "empty subject", claims: &auth.Claims{Name: td.String()}},
		{
			name: "subject", want: subject,
			claims: claims,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.claims.UserID())
		})
	}
}
