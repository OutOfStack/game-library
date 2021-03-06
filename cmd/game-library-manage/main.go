package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/OutOfStack/game-library/internal/app/game-library-manage/schema"
	"github.com/OutOfStack/game-library/internal/pkg/database"
	"github.com/OutOfStack/game-library/internal/pkg/util"
)

func main() {
	type config struct {
		DB struct {
			Host       string `mapstructure:"APP_HOST"`
			Name       string `mapstructure:"APP_NAME"`
			User       string `mapstructure:"APP_USER"`
			Password   string `mapstructure:"APP_PASSWORD"`
			RequireSSL bool   `mapstructure:"APP_REQUIRESSL"`
		} `mapstructure:",squash"`
	}

	cfg := config{}
	if err := util.LoadConfig(".", "app", "env", &cfg); err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	fmt.Printf("Host: %s, Name: %s, User: %s, RequireSSL: %v\n", cfg.DB.Host, cfg.DB.Name, cfg.DB.User, cfg.DB.RequireSSL)

	db, err := database.Open(database.Config{
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		RequireSSL: cfg.DB.RequireSSL,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db, true); err != nil {
			log.Fatalf("applying migrations %v", err)
		}
		log.Print("migration complete")
		return
	case "rollback":
		if err := schema.Migrate(db, false); err != nil {
			log.Fatalf("rollback last migration %v", err)
		}
		log.Print("migration rollback complete")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatalf("applying seeds %v", err)
		}
		log.Print("Seed data inserted")
		return
	}
}
