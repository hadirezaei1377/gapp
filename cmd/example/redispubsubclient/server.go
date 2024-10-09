package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"gapp/adapter/redis"
	"gapp/config"
	"gapp/contract/golang/matching"
	"gapp/entity"
	"gapp/pkg/slice"

	"google.golang.org/protobuf/proto"
)

func main() {
	cfg := config.Load("config.yml")
	redisAdapter := redis.New(cfg.Redis)
	topic := "matching.users_matched"
	mu := entity.MatchedUsers{
		Category: entity.FootballCategory,
		UserIDs:  []uint{1, 4},
	}
	pbMu := matching.MatchedUsers{
		Category: string(mu.Category),
		UserIds:  slice.MapFromUintToUint64(mu.UserIDs),
	}
	payload, err := proto.Marshal(&pbMu)
	if err != nil {
		panic(err)
	}
	payloadStr := base64.StdEncoding.EncodeToString(payload)
	if err := redisAdapter.Client().Publish(context.Background(), topic, payloadStr).Err(); err != nil {
		panic(fmt.Sprintf("publish err: %v", err))
	}
}
