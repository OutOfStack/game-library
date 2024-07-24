package middleware

import (
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Errors handles errors in middleware chain
func Errors(log *zap.Logger) gin.HandlerFunc {
	h := func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				log.Error(e)
			}

			err := c.Errors.Last().Err

			// reset errors as they were logged
			c.Errors = c.Errors[0:0]

			web.RespondError(c, err)
		}
	}

	return h
}
