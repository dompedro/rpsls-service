package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"rpsls/rpslsapi"
)

type ScoreboardStore struct {
	Client
}

func NewScoreboardStore(client Client) ScoreboardStore {
	return ScoreboardStore{client}
}

func (ss ScoreboardStore) Scoreboard(key string, size int64) ([]rpslsapi.RoundResults, error) {
	lastResults, err := ss.LRange(context.Background(), key, 0, size-1).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	scoreboard := make([]rpslsapi.RoundResults, len(lastResults))
	for i := range lastResults {
		var results rpslsapi.RoundResults
		err = json.Unmarshal([]byte(lastResults[i]), &results)
		if err != nil {
			return nil, err
		}
		scoreboard[i] = results
	}
	return scoreboard, nil
}

func (ss ScoreboardStore) Append(key string, size int64, results *rpslsapi.RoundResults) error {
	marshal, err := json.Marshal(results)
	if err != nil {
		return err
	}

	value := string(marshal)
	_, err = ss.LPush(context.Background(), key, value).Result()
	if err != nil {
		return err
	}

	ss.LTrim(context.Background(), key, 0, size-1)
	return nil
}

func (ss ScoreboardStore) Clear(key string) error {
	_, err := ss.Del(context.Background(), key).Result()
	return err
}
