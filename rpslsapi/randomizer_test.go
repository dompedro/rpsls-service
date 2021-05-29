package rpslsapi

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type RandomizerClientMock struct {
	mock.Mock
}

func (rcm *RandomizerClientMock) RandomNumber() (*RandomNumberResponse, error) {
	args := rcm.Called()
	return args.Get(0).(*RandomNumberResponse), args.Error(1)
}

func TestExternalRandomizerService_RandomInt(t *testing.T) {
	testCases := []struct {
		name                 string
		randomizerError      error
		randomNumberResponse *RandomNumberResponse
		expectedInt          int
		expectedError        error
	}{
		{
			name:                 "success: return generated random int",
			randomNumberResponse: &RandomNumberResponse{RandomNumber: 27},
			expectedInt:          27,
		},
		{
			name:            "failure: if randomizer returns an error, return ErrRandomNumberGenerationFailed",
			randomizerError: errors.New("unknown randomizer error"),
			expectedError:   ErrRandomNumberGenerationFailed,
		},
	}

	clientMock := RandomizerClientMock{}
	service := NewExternalRandomizerService(&clientMock)

	for _, tc := range testCases {
		clientMock.On("RandomNumber").Return(tc.randomNumberResponse, tc.randomizerError).Once()
		randomInt, err := service.RandomInt()

		if tc.expectedError != nil {
			require.NotNil(t, t, err)
			require.EqualError(t, tc.expectedError, err.Error())
		} else {
			require.Equal(t, randomInt, tc.expectedInt)
		}
	}
}
