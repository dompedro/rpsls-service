package logger

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetUp() {
	log.Logger = zerolog.New(os.Stdout).With().
		Caller().
		Stack().
		Logger()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.CallerFieldName = "loggerName"
	zerolog.ErrorFieldName = "goError"
	zerolog.ErrorStackFieldName = "stackTrace"
}

func WithReqIdAndAction(e *zerolog.Event, r *http.Request, action string) *zerolog.Event {
	return e.Str("requestId", middleware.GetReqID(r.Context())).
		Str("action", action)
}
