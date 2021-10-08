package middleware

import (
	"log"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
)

// Errors handles errors in middleware chain
func Errors(log *log.Logger) gin.HandlerFunc {

	h := func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				log.Printf("ERROR: %s\n", e)
			}

			err := c.Errors.Last().Err

			// reset errors as they were logged
			c.Errors = c.Errors[0:0]

			web.RespondError(c, err)
		}
	}

	return h
}
