package api

import (
	"io"
	"os"
	"regexp"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

func InitLogging() {
	minLevel := zerolog.InfoLevel
	if DevMode {
		minLevel = zerolog.DebugLevel
	}

	zerolog.DurationFieldUnit = time.Second
	log.Logger = log.
		Output(zerolog.ConsoleWriter{Out: redactedWriter{out: os.Stderr}, TimeFormat: time.StampMicro}).
		Level(minLevel).
		With().
		Caller().
		Logger()
}
