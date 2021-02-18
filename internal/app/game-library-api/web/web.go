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
func NewApp(logger *log.Logger) *App {
	return &App{
		router: gin.Default(),
		log:    logger,
	}
}

// Handle connect method and pattern to application handler
func (a *App) Handle(method, pattern string, fn gin.HandlerFunc) {
	a.router.Handle(method, pattern, fn)
}

// ServeHTTP serves http server
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
