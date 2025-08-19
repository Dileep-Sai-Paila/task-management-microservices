package usecase

import (
	"fmt"
	"log"
	"user_service/internal/core/users"
	"user_service/pkg/generatejwt"
	"user_service/pkg/hashpassword"
)

type userUsecase struct {
	userRepo     UserRepository
	jwtSecretKey string
}

func NewUserUsecase(repo UserRepository, jwtSecret string) UserUsecase {
	return &userUsecase{
		userRepo:     repo,
		jwtSecretKey: jwtSecret,
	}
}

func (uc *userUsecase) RegisterUser(user *users.User) error {
	hashedPassword, err := hashpassword.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}
	user.PasswordHash = hashedPassword
	err = uc.userRepo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("could not register user: %w", err)
	}
	return nil
}

func (uc *userUsecase) LoginUser(email, password string) (string, error) {
	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}
	log.Printf("Password hash from DB: %s", user.PasswordHash)
	comparisonResult := hashpassword.CheckPasswordHash(password, user.PasswordHash)
	log.Printf("Password comparison result: %v", comparisonResult)
	if !comparisonResult {
		return "", fmt.Errorf("invalid email or password")
	}
	token, err := generatejwt.GenerateToken(user, uc.jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}
	return token, nil
}

func (uc *userUsecase) GetProfile(userID int) (*users.User, error) {
	user, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve profile: %w", err)
	}
	return user, nil
}
