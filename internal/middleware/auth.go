package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/OutOfStack/game-library/internal/auth"
	"go.uber.org/zap"
)

type tokenKey int
type claimsKey int

// Context keys for authentication and authorization
var (
	tokenCtxKey  tokenKey
	claimsCtxKey claimsKey
)

// AuthClient - auth client interface
type AuthClient interface {
	ParseToken(tokenStr string) (*auth.Claims, error)
	Verify(ctx context.Context, tokenStr string) error
}

// Authenticate checks the validity of a token
func Authenticate(log *zap.Logger, authClient AuthClient) func(http.Handler) http.Handler {
	h := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			// if no Authorization header provided, return 401
			if authHeader == "" {
				web.RespondError(w, web.NewErrorFromMessage("no Authorization header found", http.StatusUnauthorized))
				return
			}

			tokenStr := auth.ExtractToken(authHeader)
			// if no Bearer token provided, return 401
			if tokenStr == "" {
				web.RespondError(w, web.NewErrorFromMessage("no Bearer token found", http.StatusUnauthorized))
				return
			}

			// if token is not valid, return 401
			err := authClient.Verify(r.Context(), tokenStr)
			if err != nil {
				log.Error("verifying token", zap.Error(err))
				statusCode := http.StatusUnauthorized
				if errors.Is(err, auth.ErrVerifyAPIUnavailable) {
					statusCode = http.StatusBadGateway
				}
				web.RespondError(w, web.NewError(err, statusCode))
				return
			}

			// store the token in the request context
			ctx := context.WithValue(r.Context(), tokenCtxKey, tokenStr)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return h
}

// Authorize checks rights to perform certain requests
func Authorize(log *zap.Logger, authClient AuthClient, requiredRole string) func(http.Handler) http.Handler {
	h := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// retrieve the token from context
			token, ok := r.Context().Value(tokenCtxKey).(string)
			// if no value in context return 500 as it is unexpected
			if !ok || token == "" {
				log.Error("no token in request context")
				web.RespondError(w, web.NewErrorFromStatusCode(http.StatusInternalServerError))
				return
			}

			// parse the token
			claims, err := authClient.ParseToken(token)
			if err != nil {
				log.Error("parsing token", zap.Error(err))
				web.RespondError(w, web.NewErrorFromStatusCode(http.StatusInternalServerError))
				return
			}

			// check user's role
			if claims.UserRole != requiredRole {
				log.Warn("access denied", zap.String("required_role", requiredRole), zap.String("got_role", claims.UserRole))
				web.RespondError(w, web.NewErrorFromMessage("access denied", http.StatusForbidden))
				return
			}

			// store claims in the request context
			ctx := context.WithValue(r.Context(), claimsCtxKey, *claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	return h
}

// GetClaims returns user claims from context
func GetClaims(ctx context.Context) (*auth.Claims, error) {
	claims, ok := ctx.Value(claimsCtxKey).(auth.Claims)
	if !ok {
		return nil, errors.New("claims not found in context")
	}
	return &claims, nil
}
