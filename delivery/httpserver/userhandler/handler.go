package userhandler

import (
	"gapp/service/authservice"

	"gapp/service/userservice"
	"gapp/validator/uservalidator"
)

type Handler struct {
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
	presenceSvc   presenceservice.Service
}

func New(authConfig authservice.Config, authSvc authservice.Service,
	userSvc userservice.Service,

	userValidator uservalidator.Validator, presenceSvc presenceservice.Service) Handler {
	return Handler{
		authConfig:    authConfig,
		authSvc:       authSvc,
		userSvc:       userSvc,
		userValidator: userValidator,
		presenceSvc:   presenceSvc,
	}
}
