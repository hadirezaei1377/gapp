package main

import (
	"context"
	"fmt"
	"gapp/adapter/redis"
	"gapp/config"
	"gapp/entity"
	"gapp/pkg/protobufencoder"
)

func main() {
	cfg := config.Load("config.yml")

	redisAdapter := redis.New(cfg.Redis)

	topic := entity.MatchingUsersMatchedEvent

	mu := entity.MatchedUsers{
		Category: entity.FootballCategory,
		UserIDs:  []uint{1, 4},
	}

	payload := protobufencoder.EncodeEvent(topic, mu)

	if err := redisAdapter.Client().Publish(context.Background(), string(topic), payload).Err(); err != nil {
		panic(fmt.Sprintf("publish err: %v", err))
	}
}
