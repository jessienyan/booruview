package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	api "github.com/jessienyan/booruview"
	"github.com/jessienyan/booruview/routes"
)

var (
	apiKeyRedact = regexp.MustCompile(`(api_key=)[a-f0-9]+`)
)

type redactedWriter struct {
	out io.Writer
}

func (rw redactedWriter) Write(data []byte) (n int, err error) {
	originalLen := len(data)
	data = apiKeyRedact.ReplaceAll(data, []byte("$1****"))

	_, err = rw.out.Write(data)
	// Need to return the original length otherwise zerolog complains about a short write
	return originalLen, err
}

func main() {
	zerolog.DurationFieldUnit = time.Second
	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: redactedWriter{out: os.Stderr}, TimeFormat: time.StampMicro}).
		With().
		Caller().
		Logger()
	timeout := 1 * time.Second

	for {
		if err := api.InitValkey(); err != nil {
			log.Err(err).Dur("retry", timeout).Msg("failed to connect to valkey")
			time.Sleep(timeout)
			timeout *= 2
		} else {
			break
		}
	}

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
