package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/game"
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

	data, err := json.Marshal(list)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		g.Log.Println("Error marshalling", err)
		return
	}
	c.Header("content-type", "application/json;charset=utf-8")
	_, err = c.Writer.Write(data)
	if err != nil {
		g.Log.Println("Error writing", err)
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

	list, err := game.Retrieve(g.DB, id)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		g.Log.Println("Error querying db", err)
		return
	}

	if list == nil {
		c.Status(http.StatusNotFound)
		g.Log.Println("Not found for id", id)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		g.Log.Println("Error marshalling", err)
		return
	}
	c.Header("content-type", "application/json;charset=utf-8")
	_, err = c.Writer.Write(data)
	if err != nil {
		g.Log.Println("Error writing", err)
	}
}
