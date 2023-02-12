package middleware

import (
	"errors"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Authenticate checks validity of token
func Authenticate(log *zap.Logger, a *auth.Auth) gin.HandlerFunc {

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
		if err := a.Verify(c.Request.Context(), tokenStr); err != nil {
			log.Error("verifying token", zap.Error(err))
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
func Authorize(log *zap.Logger, a *auth.Auth, requiredRole string) gin.HandlerFunc {

	h := func(c *gin.Context) {
		token, ok := c.Get(auth.CtxTokenKey)
		// if no value in context return 500 as it is unexpected
		if !ok {
			log.Error("no token in request context")
			c.Error(web.NewRequestError(errors.New("internal server error"), http.StatusInternalServerError))
			c.Abort()
			return
		}
		tokenStr := token.(string)
		claims, err := a.ParseToken(tokenStr)
		// if we can't parse after verification return 500 as it is unexpected
		if err != nil {
			log.Error("parsing token", zap.Error(err))
			c.Error(web.NewRequestError(errors.New("internal server error"), http.StatusInternalServerError))
			c.Abort()
			return
		}
		// if user's role is not the same as required return 403 forbidden
		if claims.UserRole != requiredRole {
			log.Warn("access denied", zap.String("expected_role", requiredRole), zap.String("got_role", claims.UserRole))
			c.Error(web.NewRequestError(errors.New("access denied"), http.StatusForbidden))
			c.Abort()
			return
		}

		c.Set(auth.CtxClaimsKey, *claims)

		c.Next()
	}

	return h
}
