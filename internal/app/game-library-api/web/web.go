package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler is a signature for all handler funcions
type Handler func(*gin.Context) error

// WrapGinHandlerFunc wraps gin.HandlerFunc and returns Handler
func WrapGinHandlerFunc(ghf gin.HandlerFunc) Handler {
	return func(c *gin.Context) error {
		ghf(c)
		return nil
	}
}

// App abstacts specific web framework
type App struct {
	router *gin.Engine
	log    *log.Logger
	mw     []Middleware
}

// NewApp constructs entrypoint for web
func NewApp(logger *log.Logger, mw ...Middleware) *App {
	r := gin.Default()
	r.Use()
	return &App{
		router: gin.Default(),
		log:    logger,
		mw:     mw,
	}
}

// Handle connect method and pattern to application handler
func (a *App) Handle(method, pattern string, h Handler) {

	h = chainMiddleware(h, a.mw)

	fn := func(c *gin.Context) {
		if err := h(c); err != nil {
			a.log.Printf("ERROR: Unhandled error %v", err)
		}
	}

	a.router.Handle(method, pattern, fn)
}

// ServeHTTP serves http server
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
