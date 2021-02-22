package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/game"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Game has handler methods for dealing with games
type Game struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List returns all games
func (g *Game) List(c *gin.Context) {
	list, err := game.List(g.DB)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		g.Log.Println("Error querying db", err)
		return
	}

	err = web.Respond(c, list, http.StatusOK)
	if err != nil {
		g.Log.Println("Error responsing", err)
		return
	}
}

// Retrieve returns a game
func (g *Game) Retrieve(c *gin.Context) {

	idparam := c.Param("id")
	id, err := strconv.ParseUint(idparam, 10, 64)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		g.Log.Println("Error parsing param url param", err)
		return
	}

	game, err := game.Retrieve(g.DB, id)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		g.Log.Println("Error querying db", err)
		return
	}

	if game == nil {
		c.Status(http.StatusNotFound)
		g.Log.Println("Not found for id", id)
		return
	}

	err = web.Respond(c, game, http.StatusOK)
	if err != nil {
		g.Log.Println("Error responsing", err)
		return
	}
}

// Create decodes JSON and creates a new Game
func (g *Game) Create(c *gin.Context) {
	var pm game.PostModel
	err := web.Decode(c, &pm)
	if err != nil {
		c.Status(http.StatusBadRequest)
		g.Log.Println("Error decoding", err)
		return
	}

	game, err := game.Create(g.DB, pm)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		g.Log.Println("Error querying db", err)
		return
	}

	err = web.Respond(c, game, http.StatusCreated)
	if err != nil {
		g.Log.Println("Error responsing", err)
		return
	}
}
