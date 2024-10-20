package main

import (
	"context"
	"fmt"
	presenceClient "gapp/adapter/presence"
	"gapp/adapter/redis"
	"gapp/config"
	"gapp/delivery/httpserver"
	"gapp/logger"
	"gapp/repository/migrator"
	"gapp/repository/mysql"
	"gapp/repository/mysql/mysqlaccesscontrol"
	"gapp/repository/mysql/mysqluser"
	"gapp/repository/redis/redismatching"
	"gapp/repository/redis/redispresence"
	"gapp/scheduler"
	"gapp/service/authorizationservice"
	"gapp/service/authservice"
	"gapp/service/backofficeuserservice"
	"gapp/service/matchingservice"
	"gapp/service/presenceservice"
	"gapp/service/userservice"
	"gapp/validator/matchingvalidator"
	"gapp/validator/uservalidator"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"

	_ "net/http/pprof"
)

const (
	JwtSignKey = ""
)

func main() {
	go func() {
		// TODO - add enabler config variable
		// curl http://localhost:8099/debug/pprof/goroutine --output goroutine.o
		//  go tool pprof -http=:8086 ./goroutine.o
		http.ListenAndServe(":8099", nil)
	}()

	// TODO - read config path from command line
	cfg := config.Load("config.yml")

	logger.Logger.Named("main").Info("config", zap.Any("config", cfg))

	// TODO - add command for migrations
	mgr := migrator.New(cfg.Mysql)
	mgr.Up()

	// TODO - add struct and add these returned items as struct field
	authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc, matchingSvc, matchingV, presenceSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc,
		matchingSvc, matchingV, presenceSvc)
	go func() {
		server.Serve()
	}()

	done := make(chan bool)
	var wg sync.WaitGroup
	go func() {
		sch := scheduler.New(cfg.Scheduler, matchingSvc)

		wg.Add(1)
		sch.Start(done, &wg)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cfg.Application.GracefulShutdownTimeout)
	defer cancel()

	if err := server.Router.Shutdown(ctxWithTimeout); err != nil {
		// TODO - replace all fmt.Print.. and std log calls with logger.Logger
		fmt.Println("http server shutdown error", err)
	}

	fmt.Println("received interrupt signal, shutting down gracefully..")
	done <- true
	time.Sleep(cfg.Application.GracefulShutdownTimeout)

	// TODO - does order of ctx.Done & wg.Wait matter?
	<-ctxWithTimeout.Done()

	wg.Wait()
}

func setupServices(cfg config.Config) (
	authservice.Service, userservice.Service, uservalidator.Validator,
	backofficeuserservice.Service, authorizationservice.Service,
	matchingservice.Service, matchingvalidator.Validator,
	presenceservice.Service,
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

	presenceRepo := redispresence.New(redisAdapter)
	presenceSvc := presenceservice.New(cfg.PresenceService, presenceRepo)

	matchingRepo := redismatching.New(redisAdapter)

	// TODO - add address to config
	presenceAdapter := presenceClient.New(":8086")

	matchingSvc := matchingservice.New(cfg.MatchingService, matchingRepo, presenceAdapter, redisAdapter)

	return authSvc, userSvc, uV, backofficeUserSvc, authorizationSvc, matchingSvc, matchingV, presenceSvc
}
