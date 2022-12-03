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

// GetIDParam returns url id param
func GetIDParam(c *gin.Context) (int32, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil || id <= 0 {
		return 0, NewRequestError(errors.New("invalid id"), http.StatusBadRequest)
	}
	return int32(id), err
}
