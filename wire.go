//+build wireinject

package rpsls_gen

import (
	"github.com/google/wire"
	"rpsls/rpslsapi"
	"rpsls/rpslsapi/http"
	"rpsls/rpslsapi/storage/neo4j"
	"rpsls/rpslsapi/storage/redis"
)

func InitServer() (http.Server, func()) {
	wire.Build(
		http.NewServer,
		http.NewRouter,
		http.NewChoiceHandler,
		http.NewRoundHandler,
		http.NewScoreboardHandler,
		http.NewRandomizerClient,
		rpslsapi.NewExternalRandomizerService,
		rpslsapi.NewChoiceService,
		rpslsapi.NewRoundService,
		rpslsapi.NewScoreboardService,
		neo4j.NewDbClient,
		neo4j.NewChoiceStore,
		neo4j.NewRoundStore,
		redis.NewClient,
		redis.NewScoreboardStore,
		wire.Bind(new(rpslsapi.ChoiceStore), new(neo4j.ChoiceStore)),
		wire.Bind(new(rpslsapi.RoundStore), new(neo4j.RoundStore)),
		wire.Bind(new(rpslsapi.ScoreboardStore), new(redis.ScoreboardStore)),
		wire.Bind(new(rpslsapi.RandomizerService), new(rpslsapi.ExternalRandomizerService)),
		wire.Bind(new(rpslsapi.RandomizerClient), new(http.RandomizerClient)))

	return http.Server{}, func() {}
}
