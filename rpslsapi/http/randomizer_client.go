package http

import (
	"encoding/json"
	"net/http"

	"rpsls/rpslsapi"
)

type RandomizerClient struct {
	url string
}

func NewRandomizerClient() RandomizerClient {
	return RandomizerClient{url: rpslsapi.Config.RandomNumberServer}
}

func (rc RandomizerClient) RandomNumber() (*rpslsapi.RandomNumberResponse, error) {
	resp, err := http.Get(rc.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var randomNumberResponse rpslsapi.RandomNumberResponse
	err = json.NewDecoder(resp.Body).Decode(&randomNumberResponse)
	if err != nil {
		return nil, err
	}

	return &randomNumberResponse, nil
}
