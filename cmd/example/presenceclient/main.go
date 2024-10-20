package main

import (
	"context"
	"fmt"
	presenceClient "gapp/adapter/presence"
	"gapp/param"
	"log"
)

func main() {
	client := presenceClient.New("localhost:8086")

	resp, err := client.GetPresence(context.Background(), param.GetPresenceRequest{UserIDs: []uint{1, 2, 4}})
	if err != nil {
		log.Fatalf("GetPresence error: %v", err)
	}

	for _, item := range resp.Items {
		fmt.Println("item", item.UserID, item.Timestamp)
	}
}
