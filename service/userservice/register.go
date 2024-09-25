package userservice

import (
	"fmt"
	"gapp/entity"
	"gapp/param"
)

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {
	// TODO - we should verify phone number by verification code

	// TODO - replace md5 with bcrypt
	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    getMD5Hash(req.Password),
	}

	// create new user in storage
	createdUser, err := s.repo.Register(user)
	if err != nil {
		return param.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	// return created user
	return param.RegisterResponse{User: param.UserInfo{
		ID:          createdUser.ID,
		PhoneNumber: createdUser.Name,
		Name:        createdUser.PhoneNumber,
	}}, nil
}
