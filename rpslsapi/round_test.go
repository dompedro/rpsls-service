package rpslsapi

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type RoundStoreMock struct {
	mock.Mock
}

type ChoiceServiceMock struct {
	mock.Mock
}

type ScoreboardServiceMock struct {
	mock.Mock
}

func (rsm *RoundStoreMock) SimulateRound(choice1ID, choice2ID int64) (*Round, error) {
	args := rsm.Called(choice1ID, choice2ID)
	return args.Get(0).(*Round), args.Error(1)
}

func (csm *ChoiceServiceMock) Choices() ([]Choice, error) {
	args := csm.Called()
	return args.Get(0).([]Choice), args.Error(1)
}

func (csm *ChoiceServiceMock) Choice(id int64) (*Choice, error) {
	args := csm.Called(id)
	return args.Get(0).(*Choice), args.Error(1)
}

func (csm *ChoiceServiceMock) RandomChoice() (*Choice, error) {
	args := csm.Called()
	return args.Get(0).(*Choice), args.Error(1)
}

func (ssm *ScoreboardServiceMock) Scoreboard(userID string) ([]RoundResults, error) {
	args := ssm.Called(userID)
	return args.Get(0).([]RoundResults), args.Error(1)
}

func (ssm *ScoreboardServiceMock) Append(userID string, results *RoundResults) error {
	args := ssm.Called(userID, results)
	return args.Error(0)
}

func (ssm *ScoreboardServiceMock) Clear(userID string) error {
	args := ssm.Called(userID)
	return args.Error(0)
}

func TestRoundService_Play(t *testing.T) {
	const winnerChoiceID = int64(1)
	const loserChoiceID = int64(2)
	const missingChoiceID = int64(3)

	testCases := []struct {
		name                   string
		playerChoiceID         int64
		randomComputerChoiceID int64
		expectedOutcome        string
		expectedError          error
	}{
		{
			name:                   "success: play round and win",
			playerChoiceID:         winnerChoiceID,
			randomComputerChoiceID: loserChoiceID,
			expectedOutcome:        string(Win),
			expectedError:          nil,
		},
		{
			name:                   "success: play round and tie",
			playerChoiceID:         loserChoiceID,
			randomComputerChoiceID: loserChoiceID,
			expectedOutcome:        string(Tie),
			expectedError:          nil,
		},
		{
			name:                   "success: play round and lose",
			playerChoiceID:         loserChoiceID,
			randomComputerChoiceID: winnerChoiceID,
			expectedOutcome:        string(Lose),
			expectedError:          nil,
		},
		{
			name:                   "failure: if choice not found, return error",
			playerChoiceID:         missingChoiceID,
			randomComputerChoiceID: winnerChoiceID,
			expectedError:          ErrChoiceNotFound,
		},
	}

	storeMock := RoundStoreMock{}
	choiceServiceMock := ChoiceServiceMock{}
	scoreboardServiceMock := ScoreboardServiceMock{}
	service := NewRoundService(&storeMock, &choiceServiceMock, &scoreboardServiceMock)
	choiceServiceMock.On("Choice", winnerChoiceID).Return(&Choice{ID: winnerChoiceID}, nil)
	choiceServiceMock.On("Choice", loserChoiceID).Return(&Choice{ID: loserChoiceID}, nil)
	choiceServiceMock.On("Choice", int64(3)).Return((*Choice)(nil), ErrChoiceNotFound)
	storeMock.On("SimulateRound", winnerChoiceID, loserChoiceID).
		Return(&Round{WinnerID: winnerChoiceID, LoserID: loserChoiceID}, nil)
	storeMock.On("SimulateRound", loserChoiceID, winnerChoiceID).
		Return(&Round{WinnerID: winnerChoiceID, LoserID: loserChoiceID}, nil)
	scoreboardServiceMock.On("Append", mock.Anything, mock.Anything).Return(nil)

	for _, tc := range testCases {
		choiceServiceMock.On("RandomChoice").Return(&Choice{ID: tc.randomComputerChoiceID}, nil).Once()

		results, err := service.Play(&RoundSettings{Player: tc.playerChoiceID})

		if tc.expectedError != nil {
			require.NotNil(t, t, err)
			require.EqualError(t, tc.expectedError, err.Error())
		} else {
			require.NotNil(t, results)
			require.Equal(t, tc.expectedOutcome, results.Results)
			choiceServiceMock.AssertCalled(t, "RandomChoice")
			choiceServiceMock.AssertCalled(t, "Choice", tc.playerChoiceID)
			if results.Results != string(Tie) {
				storeMock.AssertCalled(t, "SimulateRound", tc.playerChoiceID, tc.randomComputerChoiceID)
			}
			scoreboardServiceMock.AssertCalled(t, "Append", mock.Anything, mock.Anything)
		}
	}
}
