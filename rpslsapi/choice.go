package rpslsapi

import "errors"

var ErrChoiceNotFound = errors.New("choice not found")

type Choice struct {
	ID   int64  `db:"id" json:"id"`
	Name string `json:"name"`
}

type ChoiceService interface {
	Choices() ([]Choice, error)
	RandomChoice() (*Choice, error)
	Choice(id int64) (*Choice, error)
}

type ChoiceStore interface {
	Choices() ([]Choice, error)
	Choice(id int64) (*Choice, error)
}

type ChoiceServiceImpl struct {
	store      ChoiceStore
	randomizer RandomizerService
}

func NewChoiceService(store ChoiceStore, randomizer RandomizerService) ChoiceService {
	return ChoiceServiceImpl{store: store, randomizer: randomizer}
}

func (cs ChoiceServiceImpl) Choices() ([]Choice, error) {
	choices, err := cs.store.Choices()
	if err == nil && choices == nil {
		choices = []Choice{}
	}
	return choices, err
}

func (cs ChoiceServiceImpl) Choice(id int64) (*Choice, error) {
	return cs.store.Choice(id)
}

func (cs ChoiceServiceImpl) RandomChoice() (*Choice, error) {
	randomInt, err := cs.randomizer.RandomInt()
	if err != nil {
		return nil, err
	}
	choices, err := cs.Choices()
	if err != nil {
		return nil, err
	}

	return &choices[randomInt%len(choices)], nil
}
