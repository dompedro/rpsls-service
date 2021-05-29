package rpslsapi

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ChoiceStoreMock struct {
	mock.Mock
}

type RandomizerMock struct {
	mock.Mock
}

func (csm *ChoiceStoreMock) Choices() ([]Choice, error) {
	args := csm.Called()
	return args.Get(0).([]Choice), args.Error(1)
}

func (csm *ChoiceStoreMock) Choice(id int64) (*Choice, error) {
	args := csm.Called(id)
	return args.Get(0).(*Choice), args.Error(1)
}

func (rm *RandomizerMock) RandomInt() (int, error) {
	args := rm.Called()
	return args.Get(0).(int), args.Error(1)
}

var baseChoices = []Choice{
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

var unknownDBError = errors.New("unknown DB error")
var unknownRandomizerError = errors.New("unknown randomizer error")

func TestChoiceService_Choices(t *testing.T) {
	testCases := []struct {
		name            string
		storeError      error
		existingChoices []Choice
		expectedChoices []Choice
		expectedError   error
	}{
		{
			name:            "success: return choices list",
			existingChoices: baseChoices,
			expectedChoices: baseChoices,
			expectedError:   nil,
		},
		{
			name:            "success: if no choice exists, return empty slice",
			existingChoices: []Choice{},
			expectedChoices: []Choice{},
			expectedError:   nil,
		},
		{
			name:            "success: if store returns nil, return empty slice",
			existingChoices: nil,
			expectedChoices: []Choice{},
			expectedError:   nil,
		},
		{
			name:            "failure: if store returns unknown error, propagate it",
			storeError:      unknownDBError,
			existingChoices: baseChoices,
			expectedChoices: nil,
			expectedError:   unknownDBError,
		},
	}

	storeMock := ChoiceStoreMock{}
	service := NewChoiceService(&storeMock, nil)

	for _, tc := range testCases {
		storeMock.On("Choices").Return(tc.existingChoices, tc.storeError).Once()
		choices, err := service.Choices()

		if tc.expectedError != nil {
			require.NotNil(t, t, err)
			require.EqualError(t, tc.expectedError, err.Error())
		} else {
			require.NotNil(t, choices)
			require.Equal(t, len(tc.expectedChoices), len(choices))
			storeMock.AssertCalled(t, "Choices")

			retChoicesMap := make(map[int64]Choice, len(choices))
			for i := range choices {
				retChoice := choices[i]
				retChoicesMap[retChoice.ID] = retChoice
			}

			for i := range tc.expectedChoices {
				expChoice := tc.expectedChoices[i]
				if retChoice, found := retChoicesMap[expChoice.ID]; found {
					require.Equal(t, expChoice, retChoice)
				} else {
					require.Fail(t, "choice with ID %v not returned", expChoice.ID)
				}
			}
		}
	}
}

func TestChoiceService_Choice(t *testing.T) {
	testCases := []struct {
		name           string
		givenID        int64
		expectedChoice *Choice
		expectedError  error
	}{
		{
			name:           "success: return found choice",
			givenID:        1,
			expectedChoice: &baseChoices[0],
			expectedError:  nil,
		},
		{
			name:          "failure: if choice not found, return error",
			givenID:       2,
			expectedError: ErrChoiceNotFound,
		},
		{
			name:          "failure: if store returns unknown error, propagate it",
			givenID:       3,
			expectedError: unknownDBError,
		},
	}

	storeMock := ChoiceStoreMock{}
	service := NewChoiceService(&storeMock, nil)
	storeMock.On("Choice", int64(1)).Return(&baseChoices[0], nil)
	storeMock.On("Choice", int64(2)).Return((*Choice)(nil), ErrChoiceNotFound)
	storeMock.On("Choice", int64(3)).Return((*Choice)(nil), unknownDBError)

	for _, tc := range testCases {
		choice, err := service.Choice(tc.givenID)

		if tc.expectedError != nil {
			require.NotNil(t, t, err)
			require.EqualError(t, tc.expectedError, err.Error())
		} else {
			require.NotNil(t, choice)
			require.Equal(t, tc.expectedChoice, choice)
			storeMock.AssertCalled(t, "Choice", tc.givenID)
		}
	}
}

func TestChoiceService_RandomChoice(t *testing.T) {
	testCases := []struct {
		name            string
		randomInt       int
		storeError      error
		randomizerError error
		expectedError   error
	}{
		{
			name:          "success: return random choice",
			randomInt:     len(baseChoices) + 10,
			expectedError: nil,
		},
		{
			name:          "failure: if store returns unknown error, propagate it",
			randomInt:     0,
			storeError:    unknownDBError,
			expectedError: unknownDBError,
		},
		{
			name:            "failure: if randomizer returns unknown error, propagate it",
			randomInt:       0,
			randomizerError: unknownRandomizerError,
			expectedError:   unknownRandomizerError,
		},
	}

	storeMock := ChoiceStoreMock{}
	randomizerMock := RandomizerMock{}
	service := NewChoiceService(&storeMock, &randomizerMock)

	for _, tc := range testCases {
		storeMock.On("Choices").Return(baseChoices, tc.storeError).Once()
		randomizerMock.On("RandomInt").Return(tc.randomInt, tc.randomizerError).Once()
		choice, err := service.RandomChoice()

		if tc.expectedError != nil {
			require.NotNil(t, t, err)
			require.EqualError(t, tc.expectedError, err.Error())
		} else {
			require.NotNil(t, choice)
			randomizerMock.AssertCalled(t, "RandomInt")

			found := false
			for i := range baseChoices {
				if baseChoices[i] == *choice {
					found = true
					break
				}
			}
			if !found {
				require.Fail(t, "returned non existing random choice")
			}
		}
	}
}
