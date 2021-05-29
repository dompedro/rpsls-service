package neo4j

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"rpsls/rpslsapi"
)

const allChoicesQuery = "MATCH (c:Choice) RETURN id(c) as id, c.name as name"
const choiceByIdQuery = "MATCH (c:Choice) WHERE id(c) = $id RETURN id(c) as id, c.name as name"

type ChoiceStore struct {
	DbClient
}

func NewChoiceStore(dbClient DbClient) ChoiceStore {
	return ChoiceStore{dbClient}
}

func (cs ChoiceStore) Choices() ([]rpslsapi.Choice, error) {
	session := cs.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: cs.databaseName})
	defer CloseDBResource(session)

	choices, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		records, err := transaction.Run(allChoicesQuery, nil)
		if err != nil {
			return nil, err
		}
		var result []rpslsapi.Choice

		for records.Next() {
			record := records.Record()
			id, _ := record.Get("id")
			name, _ := record.Get("name")
			result = append(result, rpslsapi.Choice{ID: id.(int64), Name: name.(string)})
		}
		return result, nil
	})
	if err != nil {
		return nil, err
	}

	return choices.([]rpslsapi.Choice), nil
}

func (cs ChoiceStore) Choice(id int64) (*rpslsapi.Choice, error) {
	session := cs.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: cs.databaseName})
	defer CloseDBResource(session)

	choices, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(choiceByIdQuery,
			map[string]interface{}{"id": id})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			record := result.Record()
			name, _ := record.Get("name")
			return &rpslsapi.Choice{ID: id, Name: name.(string)}, nil
		}
		return nil, rpslsapi.ErrChoiceNotFound
	})
	if err != nil {
		return nil, err
	}

	return choices.(*rpslsapi.Choice), nil
}
