package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"rpsls/rpslsapi"
	"rpsls/rpslsapi/logger"
)

type RoundHandler struct {
	service rpslsapi.RoundService
}

func NewRoundHandler(roundService rpslsapi.RoundService) RoundHandler {
	return RoundHandler{service: roundService}
}

func (ch *RoundHandler) addRoutes(r chi.Router) {
	r.Post("/", ch.handlePlay)
}

func (ch *RoundHandler) handlePlay(w http.ResponseWriter, r *http.Request) {
	var settings rpslsapi.RoundSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		writeJsonResponse(ErrorResponse{Code: UnprocessableBody, Message: err.Error()},
			http.StatusUnprocessableEntity, w, r, "playRound")
		logger.WithReqIdAndAction(log.Debug().Err(err), r, "playRound").
			Msg("failed to parse request")
		return
	}

	result, err := ch.service.Play(&settings)
	if err != nil {
		if err == rpslsapi.ErrChoiceNotFound {
			writeJsonResponse(ErrorResponse{Code: EntityNotFound, Message: "choice not found"},
				http.StatusNotFound, w, r, "playRound")
			logger.WithReqIdAndAction(log.Debug(), r, "playRound").
				Int64("id", settings.Player).
				Msg("choice not found")
			return
		}

		writeJsonResponse(ErrorResponse{Code: UnknownError, Message: "failed to play match"},
			http.StatusInternalServerError, w, r, "playRound")
		logger.WithReqIdAndAction(log.Error().Stack().Err(err), r, "playRound").
			Msg("failed to play match")
		return
	}

	writeJsonResponse(result, http.StatusOK, w, r, "playRound")
}
