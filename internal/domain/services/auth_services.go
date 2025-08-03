package services

import "github.com/anggitaputra13/mcs-auth/internal/domain/entities"

type AuthService interface {
	Register(user *entities.User) error
	Login(email, password string) (*entities.User, string, error)
	Logout(token string) error
	ValidateToken(token string) (string, error)
}
