package igdb

import (
	"log"

	"github.com/OutOfStack/game-library/internal/appconf"
)

// Client represents dependencies for igdb client
type Client struct {
	log   *log.Logger
	conf  appconf.IGDB
	token *token
}

// New constructs IGDB instance
func New(log *log.Logger, conf appconf.IGDB) (*Client, error) {
	return &Client{
		log:   log,
		token: &token{},
		conf:  conf,
	}, nil
}
