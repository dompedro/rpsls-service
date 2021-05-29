package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

type Router struct {
	chi.Router
}

func NewRouter(choiceHandler ChoiceHandler, roundHandler RoundHandler, scoreboardHandler ScoreboardHandler) Router {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://codechallenge.boohma.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Route("/", choiceHandler.addRoutes)
	router.Route("/play", roundHandler.addRoutes)
	router.Route("/scoreboard", scoreboardHandler.addRoutes)

	return Router{router}
}
