package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

// User roles
const (
	RoleModerator      = "moderator"
	RolePublisher      = "publisher"
	RoleRegisteredUser = "user"
)

// Claims represents jwt claims
type Claims struct {
	jwt.RegisteredClaims
	UserRole string `json:"user_role,omitempty"`
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
}

// UserID return user id from claims
func (c *Claims) UserID() string {
	return c.Subject
}
