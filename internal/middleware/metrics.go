package middleware

import (
	"expvar"

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
func Metrics() gin.HandlerFunc {

	h := func(c *gin.Context) {
		c.Next()

		m.req.Add(1)

		if len(c.Errors) > 0 {
			m.err.Add(1)
		}
	}

	return h
}
