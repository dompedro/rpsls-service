package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"rpsls/rpslsapi"
)

type ScoreboardServiceMock struct {
	mock.Mock
}

func (ssm *ScoreboardServiceMock) Scoreboard(userID string) ([]rpslsapi.RoundResults, error) {
	args := ssm.Called(userID)
	return args.Get(0).([]rpslsapi.RoundResults), args.Error(1)
}

func (ssm *ScoreboardServiceMock) Append(userID string, results *rpslsapi.RoundResults) error {
	args := ssm.Called(userID, results)
	return args.Error(0)
}

func (ssm *ScoreboardServiceMock) Clear(userID string) error {
	args := ssm.Called(userID)
	return args.Error(0)
}

func TestGetScoreboardRequest(t *testing.T) {
	results := []rpslsapi.RoundResults{
		{
			Results:  string(rpslsapi.Tie),
			Player:   1,
			Computer: 1,
		},
		{
			Results:  string(rpslsapi.Win),
			Player:   2,
			Computer: 5,
		},
	}

	testCases := []struct {
		name               string
		resultsFromService []rpslsapi.RoundResults
		serviceError       error
		expectedResults    []rpslsapi.RoundResults
		expectedStatus     int
	}{
		{
			name:               "success: return choices",
			resultsFromService: results,
			expectedResults:    results,
			expectedStatus:     http.StatusOK,
		},
		{
			name:           "failure: if an unknown error happens, return 500",
			serviceError:   errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	serviceMock := ScoreboardServiceMock{}
	router := NewRouter(ChoiceHandler{}, RoundHandler{}, NewScoreboardHandler(&serviceMock))

	for _, tc := range testCases {
		serviceMock.On("Scoreboard", mock.Anything).Return(tc.resultsFromService, tc.serviceError).Once()

		req := httptest.NewRequest("GET", "/scoreboard", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, tc.expectedStatus, rr.Code)
		if tc.expectedResults != nil {
			var returnedBody []rpslsapi.RoundResults
			err := json.Unmarshal(rr.Body.Bytes(), &returnedBody)
			require.NoError(t, err)
			require.EqualValues(t, tc.expectedResults, returnedBody)
		}
	}
}

func TestClearScoreboardRequest(t *testing.T) {
	testCases := []struct {
		name           string
		serviceError   error
		expectedStatus int
	}{
		{
			name:           "success: return 200",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "failure: if an unknown error happens, return 500",
			serviceError:   errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	serviceMock := ScoreboardServiceMock{}
	router := NewRouter(ChoiceHandler{}, RoundHandler{}, NewScoreboardHandler(&serviceMock))

	for _, tc := range testCases {
		serviceMock.On("Clear", mock.Anything).Return(tc.serviceError).Once()

		req := httptest.NewRequest("DELETE", "/scoreboard", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		require.Equal(t, tc.expectedStatus, rr.Code)
	}
}
