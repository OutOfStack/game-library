package middleware

import (
	"log"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
)

// Errors handles errors in middleware chain
func Errors(log *log.Logger) web.Middleware {

	f := func(before web.Handler) web.Handler {

		h := func(c *gin.Context) error {
			if err := before(c); err != nil {
				log.Printf("ERROR: %v", err)

				if err := web.RespondError(c, err); err != nil {
					return err
				}
			}

			return nil
		}

		return h
	}

	return f
}
