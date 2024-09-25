package userservice

import (
	"fmt"
	"gapp/dto"
	"gapp/entity"
)

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
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
		return dto.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	// return created user
	return dto.RegisterResponse{User: dto.UserInfo{
		ID:          createdUser.ID,
		PhoneNumber: createdUser.Name,
		Name:        createdUser.PhoneNumber,
	}}, nil
}
