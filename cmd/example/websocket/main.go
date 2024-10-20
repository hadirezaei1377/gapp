package main

import (
	"fmt"
	"gapp/entity"
	"gapp/pkg/protobufencoder"
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func printDecodedNotification() {
	d := protobufencoder.EncodeNotification(entity.Notification{
		EventType: "ping",
		Payload:   "hello",
	})

	fmt.Println("d", d)
}

func main() {
	printDecodedNotification()

	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
			panic(err)
		}

		defer conn.Close()

		done := make(chan bool)

		go readMessage(conn, done)

		<-done
		//go func() {
		//	defer conn.Close()
		//
		//	for {
		//		msg, op, err := wsutil.ReadClientData(conn)
		//		if err != nil {
		//			panic(err)
		//			// handle error
		//		}
		//		err = wsutil.WriteServerMessage(conn, op, msg)
		//		if err != nil {
		//			panic(err)
		//			// handle error
		//		}
		//	}
		//}()
	}))
}

func readMessage(conn net.Conn, done chan<- bool) {
	for {
		msg, opCode, err := wsutil.ReadClientData(conn)
		if err != nil {

			log.Print(err)
			done <- true
		}

		notif := protobufencoder.DecodeNotification(string(msg))

		fmt.Println("opCode", opCode)
		fmt.Println("notif", notif)

	}
}
