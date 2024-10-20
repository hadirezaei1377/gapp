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

	subscriber := redisAdapter.Client().Subscribe(context.Background(), string(topic))

	for {
		msg, err := subscriber.ReceiveMessage(context.Background())
		if err != nil {
			panic(err)
		}

		switch entity.Event(msg.Channel) {
		case topic:
			processUsersMatchedEvent(msg.Channel, msg.Payload)
		default:
			fmt.Println("invalid topic", msg.Channel)
		}
	}
}

func processUsersMatchedEvent(topic string, data string) {
	//payload, err := base64.StdEncoding.DecodeString(data)
	//if err != nil {
	//	panic(err)
	//}
	//
	//pbMu := matching.MatchedUsers{}
	//if err := proto.Unmarshal(payload, &pbMu); err != nil {
	//	panic(err)
	//}
	//
	//mu := entity.MatchedUsers{
	//	Category: entity.Category(pbMu.Category),
	//	UserIDs:  slice.MapFromUint64ToUint(pbMu.UserIds),
	//}

	payload := protobufencoder.DecodeEvent(entity.Event(topic), data)
	mu, ok := payload.(entity.MatchedUsers)
	if !ok {
		panic(ok)
	}

	fmt.Println("Received message from " + topic + " topic.")
	fmt.Printf("matched users %+v\n", mu)
}
