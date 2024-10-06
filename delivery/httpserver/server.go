package httpserver

import (
	"fmt"
	"gapp/config"
	"gapp/delivery/httpserver/backofficeuserhandler"
	"gapp/delivery/httpserver/matchinghandler"
	"gapp/delivery/httpserver/userhandler"
	"gapp/service/authorizationservice"
	"gapp/service/authservice"
	"gapp/service/backofficeuserservice"
	"gapp/service/matchingservice"
	"gapp/service/userservice"
	"gapp/validator/matchingvalidator"
	"gapp/validator/uservalidator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config                config.Config
	userHandler           userhandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler
	matchingHandler       matchinghandler.Handler
	Router                *echo.Echo
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service,
	userValidator uservalidator.Validator,
	backofficeUserSvc backofficeuserservice.Service, authorizationSvc authorizationservice.Service,
	matchingSvc matchingservice.Service,
	matchingValidator matchingvalidator.Validator) Server {
	return Server{
		Router:                echo.New(),
		config:                config,
		userHandler:           userhandler.New(config.Auth, authSvc, userSvc, userValidator),
		backofficeUserHandler: backofficeuserhandler.New(config.Auth, authSvc, backofficeUserSvc, authorizationSvc),
	}
}

func (s Server) Serve() {

	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())
	// Routes
	s.Router.GET("/health-check", s.healthCheck)

	s.userHandler.SetRoutes(s.Router)
	s.backofficeUserHandler.SetRoutes(s.Router)
	s.matchingHandler.SetRoutes(s.Router)

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}
