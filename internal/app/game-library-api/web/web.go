package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// App abstacts specific web framework
type App struct {
	router *gin.Engine
	log    *log.Logger
}

// NewApp constructs entrypoint for web
func NewApp(logger *log.Logger, mw ...gin.HandlerFunc) *App {
	r := gin.Default()
	r.Use(mw...)
	return &App{
		router: r,
		log:    logger,
	}
}

// Handle connect method and pattern to application handler
func (a *App) Handle(method, pattern string, h gin.HandlerFunc) {
	a.router.Handle(method, pattern, h)
}

// ServeHTTP serves http server
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
