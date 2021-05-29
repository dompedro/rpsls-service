package neo4j

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"rpsls/rpslsapi"
)

const beatsRelationshipQuery = "MATCH (choice1:Choice)-[BEATS]-(choice2) " +
	"WHERE id(choice1) = $choice1ID AND id(choice2) = $choice2ID " +
	"RETURN id(startNode(BEATS)) as winnerChoiceID, id(endNode(BEATS)) as loserChoiceID, BEATS.with as action"

type RoundStore struct {
	DbClient
}

func NewRoundStore(dbClient DbClient) RoundStore {
	return RoundStore{dbClient}
}

func (cs RoundStore) SimulateRound(choice1ID, choice2ID int64) (*rpslsapi.Round, error) {
	session := cs.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: cs.databaseName})
	defer CloseDBResource(session)

	round, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		records, err := transaction.Run(beatsRelationshipQuery,
			map[string]interface{}{
				"choice1ID": choice1ID,
				"choice2ID": choice2ID,
			})
		if err != nil {
			return nil, err
		}
		var result rpslsapi.Round

		if records.Next() {
			record := records.Record()
			if records.Next() {
				return nil, errors.New("multiple BEATS associations")
			}

			winnerChoiceID, _ := record.Get("winnerChoiceID")
			loserChoiceID, _ := record.Get("loserChoiceID")
			action, _ := record.Get("action")
			result.WinnerID = winnerChoiceID.(int64)
			result.LoserID = loserChoiceID.(int64)
			result.Action = action.(string)
		}
		return &result, nil
	})
	if err != nil {
		return nil, err
	}

	return round.(*rpslsapi.Round), nil
}
