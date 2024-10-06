package main

import (
	"context"
	"fmt"
	"gapp/adapter/redis"
	"gapp/config"
	"gapp/delivery/httpserver"
	"gapp/repository/migrator"
	"gapp/repository/mysql"
	"gapp/repository/mysql/mysqlaccesscontrol"
	"gapp/repository/mysql/mysqluser"
	"gapp/repository/redis/redismatching"
	"gapp/service/authorizationservice"
	"gapp/service/authservice"
	"gapp/service/backofficeuserservice"
	"gapp/service/matchingservice"
	"gapp/service/userservice"
	"gapp/validator/matchingvalidator"
	"gapp/validator/uservalidator"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
)

const (
	JwtSignKey = ""
)

func main() {
	// TODO - read config path from command line
	cfg := config.Load("config.yml")
	fmt.Printf("cfg: %+v\n", cfg)
	// TODO - add command for migrations
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()
	// TODO - add struct and add these returned items as struct field
	authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc, matchingSvc, matchingV := setupServices(cfg)
	var httpServer *echo.Echo
	go func() {
		server := httpserver.New(cfg, authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc, matchingSvc, matchingV)
		httpServer = server.Serve()
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cfg.Application.GracefulShutdownTimeout)
	defer cancel()
	if err := httpServer.Shutdown(ctxWithTimeout); err != nil {
		fmt.Println("http server shutdown error", err)
	}
	fmt.Println("received interrupt signal, shutting down gracefully..")
	<-ctxWithTimeout.Done()
}
func setupServices(cfg config.Config) (
	authservice.Service, userservice.Service, uservalidator.Validator,
	backofficeuserservice.Service, authorizationservice.Service,
	matchingservice.Service, matchingvalidator.Validator,
) {
	authSvc := authservice.New(cfg.Auth)
	MysqlRepo := mysql.New(cfg.Mysql)
	userMysql := mysqluser.New(MysqlRepo)
	userSvc := userservice.New(authSvc, userMysql)
	backofficeUserSvc := backofficeuserservice.New()
	aclMysql := mysqlaccesscontrol.New(MysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)
	uV := uservalidator.New(userMysql)
	matchingV := matchingvalidator.New()
	redisAdapter := redis.New(cfg.Redis)
	matchingRepo := redismatching.New(redisAdapter)
	matchingSvc := matchingservice.New(cfg.MatchingService, matchingRepo)
	return authSvc, userSvc, uV, backofficeUserSvc, authorizationSvc, matchingSvc, matchingV
}
