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
	subscriber := redisAdapter.Client().Subscribe(context.Background(), topic)
	for {
		msg, err := subscriber.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}
		switch msg.Channel {
		case topic:
			processUsersMatchedEvent(msg.Channel, msg.Payload)
		default:
			fmt.Println("invalid topic", msg.Channel)
		}
	}
}

func processUsersMatchedEvent(topic string, data string) {
	payload, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}
	pbMu := matching.MatchedUsers{}
	if err := proto.Unmarshal(payload, &pbMu); err != nil {
		panic(err)
	}
	mu := entity.MatchedUsers{
		Category: entity.Category(pbMu.Category),
		UserIDs:  slice.MapFromUint64ToUint(pbMu.UserIds),
	}
	fmt.Println("Received message from " + topic + " topic.")
	fmt.Printf("matched users %+v\n", mu)
}
