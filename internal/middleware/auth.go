package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/gin-gonic/gin"
)

// Authenticate checks validity of token
func Authenticate(log *log.Logger, a *auth.Auth) gin.HandlerFunc {

	h := func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		// if no Authorization header provided return 401
		if authHeader == "" {
			c.Error(web.NewRequestError(errors.New("no Authorization header found"), http.StatusUnauthorized))
			c.Abort()
			return
		}

		tokenStr := auth.ExtractToken(authHeader)
		// if no Bearer token provided return 401
		if tokenStr == "" {
			c.Error(web.NewRequestError(errors.New("no Bearer token found"), http.StatusUnauthorized))
			c.Abort()
			return
		}

		// if token is not valid return 401
		if err := a.Verify(tokenStr); err != nil {
			log.Printf("error verifying token:\n%s\n%v\n", tokenStr, err)
			if err == auth.ErrVerifyAPIUnavailable {
				c.Error(web.NewRequestError(err, http.StatusBadGateway))
			} else {
				c.Error(web.NewRequestError(err, http.StatusUnauthorized))
			}
			c.Abort()
			return
		}

		c.Set(auth.CtxTokenKey, tokenStr)

		c.Next()
	}

	return h
}

// Authorize checks rights to perform certain request
func Authorize(log *log.Logger, a *auth.Auth, requiredRole string) gin.HandlerFunc {

	h := func(c *gin.Context) {
		token, ok := c.Get(auth.CtxTokenKey)
		// if no value in context return 500 as it is unexpected
		if !ok {
			log.Println("no Token in request context")
			c.Error(web.NewRequestError(errors.New("internal server error"), http.StatusInternalServerError))
			c.Abort()
			return
		}
		tokenStr := token.(string)
		claims, err := a.ParseToken(tokenStr)
		// if we can't parse after verification return 500 as it is unexpected
		if err != nil {
			log.Printf("error parsing token: %v\n", err)
			c.Error(web.NewRequestError(errors.New("internal server error"), http.StatusInternalServerError))
			c.Abort()
			return
		}
		// if user's role is not the same as required return 403 forbidden
		if claims.UserRole != requiredRole {
			log.Printf("unauthorized: expected role %s, got %s\n", requiredRole, claims.UserRole)
			c.Error(web.NewRequestError(errors.New("internal server error"), http.StatusForbidden))
			c.Abort()
			return
		}

		c.Set(auth.CtxClaimsKey, *claims)

		c.Next()
	}

	return h
}
