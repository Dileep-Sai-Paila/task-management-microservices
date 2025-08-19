package usecase

import "user_service/internal/core/users"

// business logic operations for users.
type UserUsecase interface {
	RegisterUser(user *users.User) error
	LoginUser(email, password string) (string, error)
	GetProfile(userID int) (*users.User, error) // Added GetProfile
}

// persistence operations for users.
type UserRepository interface {
	CreateUser(user *users.User) error
	GetUserByEmail(email string) (*users.User, error)
	GetUserByID(id int) (*users.User, error) // Added GetUserByID
}
