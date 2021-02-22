package web

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler is a signature for all handler funcions
type Handler func(*gin.Context) error

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
func (a *App) Handle(method, pattern string, h Handler) {
	fn := func(c *gin.Context) {
		err := h(c)
		if err != nil {
			a.log.Printf("ERROR: %v", err)

			err = RespondError(c, err)
			if err != nil {
				a.log.Printf("ERROR: %v", err)
			}
		}
	}

	a.router.Handle(method, pattern, fn)
}

// ServeHTTP serves http server
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}
