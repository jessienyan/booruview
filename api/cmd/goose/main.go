package main

import (
	"context"
	"flag"
	"os"
	"path"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/database/migrations"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

func main() {
	api.InitLogging()
	api.LoadDatabaseEnv()

	flags := flag.NewFlagSet("goose", flag.ExitOnError)
	flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) < 1 {
		log.Fatal().Msg("usage: goose <command> [arguments]")
	}

	// Ensure the path to the database file exists
	dirPath := path.Dir(api.DatabasePath)
	if err := os.MkdirAll(dirPath, 0644); err != nil {
		log.Fatal().Msgf("failed to create path %s: %v", dirPath, err)
	}

	if err := api.InitUserDatabase(); err != nil {
		log.Fatal().Msgf("failed to open db: %v", err)
	}
	db := api.UserDB()

	// Use the embedded filesystem for the migrations
	goose.SetBaseFS(migrations.FS)

	command := args[0]
	args = args[1:]
	if err := goose.RunContext(context.Background(), command, db, "", args...); err != nil {
		log.Fatal().Msgf("goose %v: %v", command, err)
	}
}
