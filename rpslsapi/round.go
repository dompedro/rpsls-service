package rpslsapi

import "github.com/rs/zerolog/log"

type Round struct {
	WinnerID int64
	LoserID  int64
	Action   string
}

type ResultsLabel string

const (
	Win  ResultsLabel = "win"
	Tie  ResultsLabel = "tie"
	Lose ResultsLabel = "lose"
)

type RoundSettings struct {
	Player int64 `json:"player"`
}

type RoundResults struct {
	Results  string `json:"results"`
	Player   int64  `json:"player"`
	Computer int64  `json:"computer"`
}

type RoundService interface {
	Play(settings *RoundSettings) (*RoundResults, error)
}

type RoundStore interface {
	SimulateRound(choice1ID, choice2ID int64) (*Round, error)
}

type RoundServiceImpl struct {
	roundStore        RoundStore
	choiceService     ChoiceService
	scoreboardService ScoreboardService
}

func NewRoundService(roundStore RoundStore, choiceService ChoiceService, scoreboardService ScoreboardService) RoundService {
	return RoundServiceImpl{roundStore: roundStore, choiceService: choiceService, scoreboardService: scoreboardService}
}

func (rs RoundServiceImpl) Play(settings *RoundSettings) (*RoundResults, error) {
	playerChoice, err := rs.choiceService.Choice(settings.Player)
	if err != nil {
		return nil, err
	}

	computerChoice, err := rs.choiceService.RandomChoice()
	if err != nil {
		return nil, err
	}

	result := &RoundResults{
		Player:   playerChoice.ID,
		Computer: computerChoice.ID,
	}
	if playerChoice.ID == computerChoice.ID {
		result.Results = string(Tie)
		rs.saveRoundResults(result)
		return result, nil
	}

	simulatedRound, err := rs.roundStore.SimulateRound(playerChoice.ID, computerChoice.ID)
	if err != nil {
		return nil, err
	}

	if simulatedRound.WinnerID == playerChoice.ID {
		result.Results = string(Win)
	} else {
		result.Results = string(Lose)
	}

	rs.saveRoundResults(result)
	return result, nil
}

func (rs RoundServiceImpl) saveRoundResults(results *RoundResults) {
	err := rs.scoreboardService.Append(Config.DummyUserID, results)
	if err != nil {
		log.Info().Msg("failed to save round to scoreboard")
	}
}
