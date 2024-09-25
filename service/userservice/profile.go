package userservice

import (
	"gapp/dto"
	"gapp/pkg/richerror"
)

// all request inputs for interactor/service should be sanitized.

func (s Service) Profile(req dto.ProfileRequest) (dto.ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return dto.ProfileResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	return dto.ProfileResponse{Name: user.Name}, nil
}
