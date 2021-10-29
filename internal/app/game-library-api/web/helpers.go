package web

import (
	"errors"
	"net/http"
	"strconv"

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

// GameIdParam return url id param
func GetIdParam(c *gin.Context) (int64, error) {
	idparam := c.Param("id")
	id, err := strconv.ParseInt(idparam, 10, 32)
	if err != nil || id <= 0 {
		return 0, NewRequestError(errors.New("invalid id"), http.StatusBadRequest)
	}
	return id, err
}
