package repositories

import "github.com/anggitaputra13/mcs-auth/internal/domain/entities"

type UserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByID(id string) (*entities.User, error)
}
