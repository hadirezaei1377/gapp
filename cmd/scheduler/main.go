package main

import (
	"fmt"
	"gapp/config"
	"gapp/scheduler"
	"os"
	"os/signal"
	"time"
)

const (
	JwtSignKey = ""
)

func main() {
	// TODO - read config path from command line
	cfg := config.Load("config.yml")
	fmt.Printf("cfg: %+v\n", cfg)

	done := make(chan bool)

	go func() {
		sch := scheduler.New()
		sch.Start(done)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("received interrupt signal, shutting down gracefully..")
	done <- true
	time.Sleep(cfg.Application.GracefulShutdownTimeout)
}
