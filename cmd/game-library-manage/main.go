package main

import (
	"flag"
	"log"

	"github.com/OutOfStack/game-library/internal/app/game-library-manage/schema"
	"github.com/OutOfStack/game-library/internal/pkg/database"
)

func main() {
	db, err := database.Open()
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
