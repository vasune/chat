package repository

import "chat/internal/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByUserID(userID uint) (*entity.User, error)
}
