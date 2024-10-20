package authorizationservice

import (
	"gapp/entity"
	"gapp/pkg/richerror"
)

type Repository interface {
	GetUserPermissionTitles(userID uint, role entity.Role) ([]entity.PermissionTitle, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) CheckAccess(userID uint, role entity.Role, permissions ...entity.PermissionTitle) (bool, error) {
	const op = "authorizationservice.CheckAccess"

	permissionTitles, err := s.repo.GetUserPermissionTitles(userID, role)
	if err != nil {
		return false, richerror.New(op).WithErr(err)
	}

	for _, pt := range permissionTitles {
		for _, p := range permissions {
			if p == pt {
				return true, nil
			}
		}
	}

	return false, nil
}
