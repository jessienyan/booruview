package main

import (
	"os"

	api "codeberg.org/jessienyan/booruview"
	"github.com/rs/zerolog/log"
)

func main() {
	api.InitLogging()
	api.LoadDatabaseEnv()
	_, err := os.Stat(api.DatabasePath)
	if os.IsNotExist(err) {
		log.Info().Msgf("creating new db: %s", api.DatabasePath)
		if err := api.InitUserDatabase(); err != nil {
			log.Fatal().Msgf("error creating db: %v", err)
		}
		api.UserDB().Ping()
	} else if err != nil {
		log.Fatal().Msgf("error checking for db: %v", err)
	}
}
