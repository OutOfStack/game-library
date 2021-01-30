package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/game"
	"github.com/jmoiron/sqlx"
)

// Game has handler methods for dealing with games
type Game struct {
	DB *sqlx.DB
}

// List returns all games
func (g *Game) List(w http.ResponseWriter, r *http.Request) {
	list, err := game.List(g.DB)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error querying db", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marshalling", err)
		return
	}
	w.Header().Set("content-type", "application/json;charset=utf-8")
	_, err = w.Write(data)
	if err != nil {
		log.Println("Error writing", err)
	}
}
