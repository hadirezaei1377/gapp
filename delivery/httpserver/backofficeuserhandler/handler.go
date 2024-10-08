package backofficeuserhandler

import (
	"gapp/service/authorizationservice"
	"gapp/service/authservice"
	"gapp/service/backofficeuserservice"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	authorizationSvc  authorizationservice.Service
	backofficeUserSvc backofficeuserservice.Service
}

func New(authConfig authservice.Config, authSvc authservice.Service,

	backofficeUserSvc backofficeuserservice.Service, authorizationSvc authorizationservice.Service) Handler {

	return Handler{
		authConfig:        authConfig,
		authSvc:           authSvc,
		backofficeUserSvc: backofficeUserSvc,
		authorizationSvc:  authorizationSvc,
	}
}
