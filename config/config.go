package config

import (
	"gapp/adapter/redis"
	"gapp/repository/mysql"
	"gapp/scheduler"
	"gapp/service/authservice"
	"gapp/service/matchingservice"
	"gapp/service/presenceservice"
	"time"
)

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	Application     Application            `koanf:"application"`
	HTTPServer      HTTPServer             `koanf:"http_server"`
	Auth            authservice.Config     `koanf:"auth"`
	Mysql           mysql.Config           `koanf:"mysql"`
	MatchingService matchingservice.Config `koanf:"matching_service"`
	Redis           redis.Config           `koanf:"redis"`
	PresenceService presenceservice.Config `koanf:"presence_service"`
	Scheduler       scheduler.Config       `koanf:"scheduler"`
}
