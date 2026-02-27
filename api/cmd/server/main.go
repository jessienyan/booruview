package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
	"codeberg.org/jessienyan/booruview/routes"
)

func main() {
	api.InitLogging()
	timeout := 1 * time.Second

	api.LoadEnv()

	for {
		if err := api.InitValkey(); err != nil {
			log.Err(err).Dur("retry", timeout).Msg("failed to connect to valkey")
			time.Sleep(timeout)
			timeout *= 2
		} else {
			break
		}
	}

	if err := api.InitUserDatabase(); err != nil {
		log.Err(err).Msg("failed to open sqlite database")
		os.Exit(1)
	}

	gelbooru.AddRatingTagsToValkey()

	// Periodically check if the CDN hosts changed. Runs immediately on server start
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		for {
			gelbooru.UpdateCDNHosts()
			<-ticker.C
		}
	}()

	listenAddr := ":8000"
	router := routes.NewRouter()
	srv := &http.Server{
		Addr:         listenAddr,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 10,
		Handler:      router,
	}

	shuttingDown := false

	go func() {
		log.Info().Str("listen", listenAddr).Msg("server started")
		if err := srv.ListenAndServe(); err != nil && !shuttingDown {
			log.Err(err).Msg("ListenAndServe error")
			log.Info().Msg("goodbye")
			os.Exit(1)
		}
	}()

	// Graceful shutdown: https://github.com/gorilla/mux?tab=readme-ov-file#graceful-shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	shuttingDown = true

	log.Info().Msg("starting graceful shutdown")
	srv.Shutdown(ctx)
	log.Info().Msg("goodbye")
}
