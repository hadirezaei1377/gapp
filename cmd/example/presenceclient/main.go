package main

import (
	"context"
	"fmt"
	presenceClient "gapp/adapter/presence"
	"gapp/param"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8086", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := presenceClient.New(conn)
	resp, err := client.GetPresence(context.Background(), param.GetPresenceRequest{UserIDs: []uint{1, 2, 4}})
	if err != nil {
		panic(err)
	}
	for _, item := range resp.Items {
		fmt.Println("item", item.UserID, item.Timestamp)
	}
}
