package web

import (
	"errors"

	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/gin-gonic/gin"
)

// GetClaims returns user claims from context
func GetClaims(c *gin.Context) (*auth.Claims, error) {
	claimsValue, ok := c.Get(auth.CtxClaimsKey)
	if !ok {
		return nil, errors.New("claims not found in context")
	}
	claims, ok := claimsValue.(auth.Claims)
	if !ok {
		return nil, errors.New("cannot convert claims value to claims")
	}
	return &claims, nil
}
