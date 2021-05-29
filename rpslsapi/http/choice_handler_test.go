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

type ChoiceServiceMock struct {
	mock.Mock
}

func (csm *ChoiceServiceMock) Choices() ([]rpslsapi.Choice, error) {
	args := csm.Called()
	return args.Get(0).([]rpslsapi.Choice), args.Error(1)
}

func (csm *ChoiceServiceMock) Choice(id int64) (*rpslsapi.Choice, error) {
	args := csm.Called(id)
	return args.Get(0).(*rpslsapi.Choice), args.Error(1)
}

func (csm *ChoiceServiceMock) RandomChoice() (*rpslsapi.Choice, error) {
	args := csm.Called()
	return args.Get(0).(*rpslsapi.Choice), args.Error(1)
}

var baseChoices = []rpslsapi.Choice{
	{
		ID:   1,
		Name: "rock",
	},
	{
		ID:   2,
		Name: "paper",
	},
	{
		ID:   3,
		Name: "scissors",
	},
	{
		ID:   4,
		Name: "lizard",
	},
	{
		ID:   5,
		Name: "spock",
	},
}

func TestChoiceListRequest(t *testing.T) {
	testCases := []struct {
		name               string
		choicesFromService []rpslsapi.Choice
		serviceError       error
		expectedChoices    []rpslsapi.Choice
		expectedStatus     int
	}{
		{
			name:               "success: return choices",
			choicesFromService: baseChoices,
			expectedChoices:    baseChoices,
			expectedStatus:     http.StatusOK,
		},
		{
			name:               "success: if service returns nil, return empty choices slice",
			choicesFromService: nil,
			expectedChoices:    []rpslsapi.Choice{},
			expectedStatus:     http.StatusOK,
		},
		{
			name:           "failure: if an unknown error happens, return 500",
			serviceError:   errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	serviceMock := ChoiceServiceMock{}
	router := NewRouter(NewChoiceHandler(&serviceMock), RoundHandler{}, ScoreboardHandler{})

	for _, tc := range testCases {
		serviceMock.On("Choices").Return(tc.choicesFromService, tc.serviceError).Once()

		req := httptest.NewRequest("GET", "/choices", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		require.Equal(t, tc.expectedStatus, rr.Code)
		if tc.expectedChoices != nil {
			var returnedBody []rpslsapi.Choice
			err := json.Unmarshal(rr.Body.Bytes(), &returnedBody)
			require.NoError(t, err)
			require.EqualValues(t, tc.expectedChoices, returnedBody)
		}
	}
}

func TestRandomChoiceRequest(t *testing.T) {
	testCases := []struct {
		name              string
		choiceFromService *rpslsapi.Choice
		serviceError      error
		bodyExpected      bool
		expectedStatus    int
	}{
		{
			name:              "success: return random choice",
			choiceFromService: &rpslsapi.Choice{ID: 1, Name: "rock"},
			expectedStatus:    http.StatusOK,
		},
		{
			name:           "failure: if an unknown error happens, return 500",
			serviceError:   errors.New("unknown error"),
			bodyExpected:   false,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	serviceMock := ChoiceServiceMock{}
	router := NewRouter(NewChoiceHandler(&serviceMock), RoundHandler{}, ScoreboardHandler{})

	for _, tc := range testCases {
		serviceMock.On("RandomChoice").Return(tc.choiceFromService, tc.serviceError).Once()

		req := httptest.NewRequest("GET", "/choice", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		require.Equal(t, tc.expectedStatus, rr.Code)
		if tc.bodyExpected {
			var returnedBody *rpslsapi.Choice
			err := json.Unmarshal(rr.Body.Bytes(), &returnedBody)
			require.NoError(t, err)
			require.NotNil(t, returnedBody)
			require.NotZero(t, returnedBody.ID)
			require.NotZero(t, returnedBody.Name)
		}
	}
}
