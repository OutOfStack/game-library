package middleware

import (
	"expvar"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
)

var m = struct {
	req *expvar.Int
	err *expvar.Int
}{
	req: expvar.NewInt("requests"),
	err: expvar.NewInt("errors"),
}

// Metrics updates program counters
func Metrics() web.Middleware {

	f := func(before web.Handler) web.Handler {

		h := func(c *gin.Context) error {
			err := before(c)
			m.req.Add(1)

			if err != nil {
				m.err.Add(1)
			}

			return err
		}

		return h
	}

	return f
}
