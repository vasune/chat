package usecases

import (
	"chat/internal/auth"
	"chat/internal/entity"
	"chat/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(repo repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{userRepo: repo}
}

func (uc *AuthUseCase) SignUp(username, password string) (string, error) {
	existingUser, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return "", err
	}
	if existingUser != nil {
		return "", errors.New("user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := &entity.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	if err := uc.userRepo.Create(user); err != nil {
		return "", err
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *AuthUseCase) SignIn(username, password string) (string, error) {
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
