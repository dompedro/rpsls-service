package rpslsapi

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ScoreboardStoreMock struct {
	mock.Mock
}

func (ssm *ScoreboardStoreMock) Scoreboard(key string, size int64) ([]RoundResults, error) {
	args := ssm.Called(key, size)
	return args.Get(0).([]RoundResults), args.Error(1)
}

func (ssm *ScoreboardStoreMock) Append(key string, size int64, results *RoundResults) error {
	args := ssm.Called(key, size, results)
	return args.Error(0)
}

func (ssm *ScoreboardStoreMock) Clear(key string) error {
	args := ssm.Called(key)
	return args.Error(0)
}

func TestScoreboardServiceImpl_Scoreboard(t *testing.T) {
	scoreboardMockError := errors.New("store error")
	results := []RoundResults{
		{
			Results:  string(Win),
			Player:   1,
			Computer: 3,
		},
		{
			Results:  string(Lose),
			Player:   4,
			Computer: 2,
		},
	}
	testCases := []struct {
		name             string
		resultsFromStore []RoundResults
		storeError       error
		expectedResults  []RoundResults
		expectedError    error
	}{
		{
			name:             "success: call store and return results",
			resultsFromStore: results,
			expectedResults:  results,
		},
		{
			name:          "failure: if store returns unknown error, propagate it",
			storeError:    scoreboardMockError,
			expectedError: scoreboardMockError,
		},
	}

	storeMock := ScoreboardStoreMock{}
	service := NewScoreboardService(&storeMock)

	for _, tc := range testCases {
		storeMock.On("Scoreboard", mock.Anything, mock.Anything).Return(tc.resultsFromStore, tc.storeError).Once()

		scoreboard, err := service.Scoreboard("dummyUserID")

		storeMock.AssertCalled(t, "Scoreboard", mock.Anything, mock.Anything)
		if tc.expectedError != nil {
			require.NotNil(t, t, err)
			require.EqualError(t, tc.expectedError, err.Error())
		} else {
			require.NoError(t, err)
			require.EqualValues(t, tc.expectedResults, scoreboard)
		}
	}
}

func TestScoreboardServiceImpl_Append(t *testing.T) {
	scoreboardMockError := errors.New("store error")

	testCases := []struct {
		name          string
		storeError    error
		expectedError error
	}{
		{
			name: "success: call store and return nil",
		},
		{
			name:          "failure: if store returns unknown error, propagate it",
			storeError:    scoreboardMockError,
			expectedError: scoreboardMockError,
		},
	}

	storeMock := ScoreboardStoreMock{}
	service := NewScoreboardService(&storeMock)

	for _, tc := range testCases {
		storeMock.On("Append", mock.Anything, mock.Anything, mock.Anything).Return(tc.storeError).Once()

		err := service.Append("dummyUserID", nil)

		storeMock.AssertCalled(t, "Append", mock.Anything, mock.Anything, mock.Anything)
		if tc.expectedError != nil {
			require.NotNil(t, t, err)
			require.EqualError(t, tc.expectedError, err.Error())
		} else {
			require.NoError(t, err)
		}
	}
}

func TestScoreboardServiceImpl_Clear(t *testing.T) {
	scoreboardMockError := errors.New("store error")

	testCases := []struct {
		name          string
		storeError    error
		expectedError error
	}{
		{
			name: "success: call store and return nil",
		},
		{
			name:          "failure: if store returns unknown error, propagate it",
			storeError:    scoreboardMockError,
			expectedError: scoreboardMockError,
		},
	}

	storeMock := ScoreboardStoreMock{}
	service := NewScoreboardService(&storeMock)

	for _, tc := range testCases {
		storeMock.On("Clear", mock.Anything).Return(tc.storeError).Once()

		err := service.Clear("dummyUserID")

		storeMock.AssertCalled(t, "Clear", mock.Anything)
		if tc.expectedError != nil {
			require.NotNil(t, t, err)
			require.EqualError(t, tc.expectedError, err.Error())
		} else {
			require.NoError(t, err)
		}
	}
}
