// `goose` CLI but with the application migrations embedded.
// This allows us to use .go migration files for complex changes.
//
// Docs:
// - https://github.com/pressly/goose?tab=readme-ov-file#embedded-sql-migrations
// - https://github.com/pressly/goose/tree/main/examples/go-migrations
// - https://github.com/pressly/goose/blob/main/examples/go-migrations/main.go

package main

import (
	"context"
	"flag"
	"os"
	"path"

	"codeberg.org/jessienyan/booruview/database/migrations"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"

	_ "modernc.org/sqlite"
)

func main() {
	flags := flag.NewFlagSet("goose", flag.ExitOnError)
	flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) == 0 {
		log.Fatal().Msg("usage: goose <command> [arguments]")
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		log.Fatal().Msg("DATABASE_PATH must be set")
	}

	// Ensure the path to the database file exists
	dirPath := path.Dir(dbPath)
	if err := os.MkdirAll(dirPath, 0644); err != nil {
		log.Fatal().Msgf("failed to create path %s: %v", dirPath, err)
	}

	db, err := goose.OpenDBWithDriver("sqlite", dbPath)
	if err != nil {
		log.Fatal().Msgf("goose: failed to open DB: %v", err)
	}
	defer func() {
		db.Close()
	}()

	// Use the embedded filesystem for the migrations
	goose.SetBaseFS(migrations.FS)

	command := args[0]
	args = args[1:]
	if err := goose.RunContext(context.Background(), command, db, ".", args...); err != nil {
		log.Fatal().Msgf("goose %v: %v", command, err)
	}
}
