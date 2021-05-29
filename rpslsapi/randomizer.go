package rpslsapi

import (
	"errors"
)

var ErrRandomNumberGenerationFailed = errors.New("failed to generate a random number")

type RandomizerService interface {
	// RandomInt generates a random int [1, 100]
	RandomInt() (int, error)
}

type RandomNumberResponse struct {
	RandomNumber int `json:"random_number"`
}

type RandomizerClient interface {
	RandomNumber() (*RandomNumberResponse, error)
}

type ExternalRandomizerService struct {
	client RandomizerClient
}

func NewExternalRandomizerService(client RandomizerClient) ExternalRandomizerService {
	return ExternalRandomizerService{client}
}

func (ers ExternalRandomizerService) RandomInt() (int, error) {
	randomNumberResponse, err := ers.client.RandomNumber()
	if err != nil {
		return 0, ErrRandomNumberGenerationFailed
	}

	return randomNumberResponse.RandomNumber, nil
}
