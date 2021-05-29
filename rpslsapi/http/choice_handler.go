package http

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"rpsls/rpslsapi"
	"rpsls/rpslsapi/logger"
)

type ChoiceHandler struct {
	service rpslsapi.ChoiceService
}

func NewChoiceHandler(choiceService rpslsapi.ChoiceService) ChoiceHandler {
	return ChoiceHandler{service: choiceService}
}

func (ch *ChoiceHandler) addRoutes(r chi.Router) {
	r.Get("/choices", ch.handleList)
	r.Get("/choice", ch.handleRandom)
}

func (ch *ChoiceHandler) handleList(w http.ResponseWriter, r *http.Request) {
	choices, err := ch.service.Choices()
	if err != nil {
		writeJsonResponse(ErrorResponse{Code: UnknownError, Message: "list choices failed"},
			http.StatusInternalServerError, w, r, "listChoices")
		logger.WithReqIdAndAction(log.Error().Stack().Err(err), r, "listChoices").
			Msg("list choices failed")
		return
	}

	if choices == nil {
		choices = []rpslsapi.Choice{}
	}
	writeJsonResponse(choices, http.StatusOK, w, r, "listChoices")
}

func (ch *ChoiceHandler) handleRandom(w http.ResponseWriter, r *http.Request) {
	randomChoice, err := ch.service.RandomChoice()
	if err != nil {
		writeJsonResponse(ErrorResponse{Code: UnknownError, Message: "random choice failed"},
			http.StatusInternalServerError, w, r, "randomChoice")
		logger.WithReqIdAndAction(log.Error().Stack().Err(err), r, "randomChoice").
			Msg("random choice failed")
		return
	}

	writeJsonResponse(randomChoice, http.StatusOK, w, r, "randomChoice")
}
