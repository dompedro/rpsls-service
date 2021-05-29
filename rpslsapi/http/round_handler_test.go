package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"rpsls/rpslsapi"
)

type RoundServiceMock struct {
	mock.Mock
}

func (rsm *RoundServiceMock) Play(settings *rpslsapi.RoundSettings) (*rpslsapi.RoundResults, error) {
	args := rsm.Called()
	return args.Get(0).(*rpslsapi.RoundResults), args.Error(1)
}

func TestPlayRequest(t *testing.T) {
	const existingChoiceID = int64(1)
	const missingChoiceID = int64(2)

	results := &rpslsapi.RoundResults{
		Results:  string(rpslsapi.Win),
		Player:   existingChoiceID,
		Computer: 3,
	}

	testCases := []struct {
		name               string
		requestBody        []byte
		resultsFromService *rpslsapi.RoundResults
		serviceError       error
		expectedResults    *rpslsapi.RoundResults
		expectedStatus     int
	}{
		{
			name:               "success: return results",
			requestBody:        playRequestBody(existingChoiceID),
			resultsFromService: results,
			expectedResults:    results,
			expectedStatus:     http.StatusOK,
		},
		{
			name:           "failure: if a missing choice ID is sent, return 404",
			requestBody:    playRequestBody(missingChoiceID),
			serviceError:   rpslsapi.ErrChoiceNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "failure: if an unknown error happens, return 500",
			requestBody:    playRequestBody(existingChoiceID),
			serviceError:   errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "failure: if a bad body is sent, return 422",
			requestBody:    []byte("{\"player\": 123"),
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	serviceMock := RoundServiceMock{}
	router := NewRouter(ChoiceHandler{}, NewRoundHandler(&serviceMock), ScoreboardHandler{})

	for _, tc := range testCases {
		serviceMock.On("Play").Return(tc.resultsFromService, tc.serviceError).Once()

		req := httptest.NewRequest("POST", "/play", bytes.NewBuffer(tc.requestBody))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		require.Equal(t, tc.expectedStatus, rr.Code)
		if tc.expectedResults != nil {
			var returnedBody *rpslsapi.RoundResults
			err := json.Unmarshal(rr.Body.Bytes(), &returnedBody)
			require.NoError(t, err)
			require.EqualValues(t, tc.expectedResults, returnedBody)
		}
	}
}

func playRequestBody(choiceID int64) []byte {
	roundSettings := rpslsapi.RoundSettings{Player: choiceID}
	body, _ := json.Marshal(roundSettings)
	return body
}
