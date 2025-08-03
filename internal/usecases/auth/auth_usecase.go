package auth_usecase

import (
	"errors"

	"github.com/anggitaputra13/mcs-auth/internal/domain/entities"
	"github.com/anggitaputra13/mcs-auth/internal/domain/repositories"
	"github.com/anggitaputra13/mcs-auth/pkg/auth"
)

type AuthUseCase struct {
	userRepo repositories.UserRepository
	jwt      *auth.JWT
}

func NewAuthUseCase(userRepo repositories.UserRepository, jwt *auth.JWT) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo, jwt: jwt}
}

func (uc *AuthUseCase) Register(user *entities.User) error {
	existingUser, err := uc.userRepo.FindByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	return uc.userRepo.Create(user)
}

func (uc *AuthUseCase) Login(email, password string) (*entities.User, string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := user.ComparePassword(password); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := uc.jwt.GenerateToken(user.ID.String())
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (uc *AuthUseCase) Logout(token string) error {
	return uc.jwt.BlacklistToken(token)
}

func (uc *AuthUseCase) ValidateToken(token string) (string, error) {
	return uc.jwt.ValidateToken(token)
}
