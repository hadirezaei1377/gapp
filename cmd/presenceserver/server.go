package main

import (
	"gapp/adapter/redis"
	"gapp/config"
	"gapp/delivery/grpcserver/presenceserver"
	"gapp/repository/redis/redispresence"
	"gapp/service/presenceservice"
)

func main() {
	// TODO - read config path from command line
	cfg := config.Load("config.yml")

	redisAdapter := redis.New(cfg.Redis)

	presenceRepo := redispresence.New(redisAdapter)
	presenceSvc := presenceservice.New(cfg.PresenceService, presenceRepo)

	server := presenceserver.New(presenceSvc)
	server.Start()
}
