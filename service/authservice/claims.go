package authservice

import (
	"gapp/entity"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint        `json:"user_id"`
	Role   entity.Role `json:"role"`
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
