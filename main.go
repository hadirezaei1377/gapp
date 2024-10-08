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
	"os"
	"os/signal"
	"sync"
	"time"

	"google.golang.org/grpc"
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

	presenceGrpcConn, err := grpc.Dial(":8086", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer presenceGrpcConn.Close()

	// TODO - add struct and add these returned items as struct field
	authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc, matchingSvc, matchingV, presenceSvc := setupServices(cfg, presenceGrpcConn)

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
		fmt.Println("http server shutdown error", err)
	}

	fmt.Println("received interrupt signal, shutting down gracefully..")
	done <- true
	time.Sleep(cfg.Application.GracefulShutdownTimeout)

	// TODO - does order of ctx.Done & wg.Wait matter?
	<-ctxWithTimeout.Done()

	wg.Wait()
}

func setupServices(cfg config.Config, presenceGrpcConn *grpc.ClientConn) (
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

	defer conn.Close()

	presenceAdapter := presenceClient.New(presenceGrpcConn)

	matchingSvc := matchingservice.New(cfg.MatchingService, matchingRepo, presenceAdapter)

	return authSvc, userSvc, uV, backofficeUserSvc, authorizationSvc, matchingSvc, matchingV, presenceSvc
}
