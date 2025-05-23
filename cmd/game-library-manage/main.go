package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/OutOfStack/game-library/internal/app/game-library-manage/schema"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/pkg/database"
)

func main() {
	cfg, err := appconf.Get()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	db, err := database.New(ctx, cfg.GetDB().DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err = schema.Migrate(cfg.DB.DSN, true); err != nil {
			log.Printf("Apply migrations failed: %v", err)
			return
		}
		log.Print("Migration complete")
	case "rollback":
		if err = schema.Migrate(cfg.DB.DSN, false); err != nil {
			log.Printf("Rollback last migration failed: %v", err)
			return
		}
		log.Print("Migration rollback complete")
	case "seed":
		if err = schema.Seed(db); err != nil {
			log.Printf("Aapply seeds failed: %v", err)
			return
		}
		log.Print("Seed data inserted")
	default:
		fmt.Println("Unknown command, available commands:")
		fmt.Println("migrate: applies all migrations to database")
		fmt.Println("rollback: roll backs one last migration of database")
		fmt.Println("seed: applies seed data (games) to database")
	}
}
