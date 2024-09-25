package main

import (
	"fmt"
	"gapp/config"
	"gapp/delivery/httpserver"
	"gapp/repository/mysql"
	"gapp/service/authservice"
	"gapp/service/userservice"
	"gapp/validator/uservalidator"
	"time"
)

const (
	JwtSignKey                 = "jwt_secret"
	AccessTokenSubject         = "ac"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func main() {
	// TODO - read config path from command line
	cfg2 := config.Load("config.yml")
	fmt.Printf("cfg2: %+v\n", cfg2)

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8088},
		Auth: authservice.Config{
			SignKey:               JwtSignKey,
			AccessExpirationTime:  AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
			AccessSubject:         AccessTokenSubject,
			RefreshSubject:        RefreshTokenSubject,
		},
		Mysql: mysql.Config{
			Username: "gameapp",
			Password: "gameappt0lk2o20",
			Port:     3308,
			Host:     "localhost",
			DBName:   "gameapp_db",
		},
	}

	// TODO - add command for migrations
	//mgr := migrator.New(cfg.Mysql)
	//mgr.Up()

	authSvc, userSvc, userValidator := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator)

	fmt.Println("start echo server")
	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)

	MysqlRepo := mysql.New(cfg.Mysql)
	userSvc := userservice.New(authSvc, MysqlRepo)

	uV := uservalidator.New(MysqlRepo)

	return authSvc, userSvc, uV
}
