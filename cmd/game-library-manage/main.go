package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/OutOfStack/game-library/internal/app/game-library-manage/schema"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/pkg/database"
)

func main() {
	var fromFile bool
	flag.BoolVar(&fromFile, "from-file", false, "read dsn from config file instead of environment variable")
	flag.Parse()

	var dsn string
	if fromFile {
		cfg, err := appconf.Get()
		if err != nil {
			log.Fatal("read config file:", err)
		}
		dsn = cfg.GetDB().DSN
		if dsn == "" {
			log.Fatal("DB_DSN not found in config file")
		}
	} else {
		dsn = os.Getenv("DB_DSN")
		if dsn == "" {
			log.Fatal("DB_DSN environment variable is required")
		}
	}

	ctx := context.Background()
	db, err := database.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch flag.Arg(0) {
	case "migrate":
		if err = schema.Migrate(dsn, true); err != nil {
			log.Printf("Apply migrations failed: %v", err)
			return
		}
		log.Print("Migration complete")
	case "rollback":
		if err = schema.Migrate(dsn, false); err != nil {
			log.Printf("Rollback last migration failed: %v", err)
			return
		}
		log.Print("Migration rollback complete")
	case "seed":
		if err = schema.Seed(db); err != nil {
			log.Printf("Apply seeds failed: %v", err)
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
