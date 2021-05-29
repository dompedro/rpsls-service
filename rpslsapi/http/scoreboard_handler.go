package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"rpsls/rpslsapi"
	"rpsls/rpslsapi/logger"
)

type ScoreboardHandler struct {
	service rpslsapi.ScoreboardService
}

func NewScoreboardHandler(scoreboardService rpslsapi.ScoreboardService) ScoreboardHandler {
	return ScoreboardHandler{service: scoreboardService}
}

func (sh *ScoreboardHandler) addRoutes(r chi.Router) {
	r.Get("/", sh.handleScoreboard)
	r.Delete("/", sh.handleClear)
}

func (sh *ScoreboardHandler) handleScoreboard(w http.ResponseWriter, r *http.Request) {
	result, err := sh.service.Scoreboard(rpslsapi.Config.DummyUserID)
	if err != nil {
		writeJsonResponse(ErrorResponse{Code: UnknownError, Message: "failed to get scoreboard"},
			http.StatusInternalServerError, w, r, "getScoreboard")
		logger.WithReqIdAndAction(log.Error().Stack().Err(err), r, "getScoreboard").
			Msg("failed to get scoreboard")
		return
	}

	writeJsonResponse(result, http.StatusOK, w, r, "getScoreboard")
}

func (sh *ScoreboardHandler) handleClear(w http.ResponseWriter, r *http.Request) {
	err := sh.service.Clear(rpslsapi.Config.DummyUserID)
	if err != nil {
		writeJsonResponse(ErrorResponse{Code: UnknownError, Message: "failed to get scoreboard"},
			http.StatusInternalServerError, w, r, "getScoreboard")
		logger.WithReqIdAndAction(log.Error().Stack().Err(err), r, "getScoreboard").
			Msg("failed to get scoreboard")
		return
	}

	w.WriteHeader(http.StatusOK)
}
