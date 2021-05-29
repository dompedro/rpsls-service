package rpslsapi

import (
	"fmt"
)

// first placeholder is for the userID
const scoreboardKeyTemplate = "rpsls-scoreboard:%s"

type ScoreboardService interface {
	Scoreboard(userID string) ([]RoundResults, error)
	Append(userID string, results *RoundResults) error
	Clear(userID string) error
}

type ScoreboardStore interface {
	Scoreboard(key string, size int64) ([]RoundResults, error)
	Append(key string, size int64, results *RoundResults) error
	Clear(key string) error
}

type ScoreboardServiceImpl struct {
	scoreboardStore ScoreboardStore
	boardSize       int
}

func NewScoreboardService(scoreboardStore ScoreboardStore) ScoreboardService {
	return ScoreboardServiceImpl{scoreboardStore, Config.ScoreboardSize}
}

func (ss ScoreboardServiceImpl) Scoreboard(userID string) ([]RoundResults, error) {
	return ss.scoreboardStore.Scoreboard(fmt.Sprintf(scoreboardKeyTemplate, userID), int64(ss.boardSize))
}

func (ss ScoreboardServiceImpl) Append(userID string, results *RoundResults) error {
	return ss.scoreboardStore.Append(fmt.Sprintf(scoreboardKeyTemplate, userID), int64(ss.boardSize), results)
}

func (ss ScoreboardServiceImpl) Clear(userID string) error {
	return ss.scoreboardStore.Clear(fmt.Sprintf(scoreboardKeyTemplate, userID))
}
