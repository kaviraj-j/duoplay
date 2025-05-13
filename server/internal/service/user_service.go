package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kaviraj-j/duoplay/internal/model"
	"github.com/kaviraj-j/duoplay/internal/repository"
)

// UserService defines the interface for user related business logic
// type UserService interface {
// 	RegisterUser(ctx context.Context, name string) (*model.User, string, error)
// 	GetUserByID(ctx context.Context, id string) (*model.User, error)
// 	ValidateToken(ctx context.Context, tokenString string) (*model.User, error)
// }

// userService implements UserService
type UserService struct {
	userRepository repository.UserRepository
	jwtSecret      []byte
	jwtExpires     time.Duration
}

func CreateUserService(userRepository repository.UserRepository, jwtSecret []byte) (*UserService, error) {
	return &UserService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
		jwtExpires:     time.Hour * 24,
	}, nil
}

// RegisterUser registers a new user in user repo
func (service *UserService) RegisterUser(ctx context.Context, name string) (*model.User, string, error) {
	user := &model.User{
		ID:   uuid.New().String(),
		Name: name,
	}
	err := service.userRepository.Create(ctx, user)
	if err != nil {
		return nil, "", err
	}
	token, err := service.generateJWT(user)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func (service *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return service.userRepository.FindByID(ctx, id)
}

func (service *UserService) ValidateToken(ctx context.Context, tokenString string) (*model.User, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return service.jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	// extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// get user ID from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid user claim")
	}

	// find user by ID
	return service.userRepository.FindByID(ctx, userID)
}

// generateJWT creates a new JWT token for a user
func (service *UserService) generateJWT(user *model.User) (string, error) {
	// Create claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"name":    user.Name,
		"exp":     time.Now().Add(service.jwtExpires).Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return the token
	return token.SignedString(service.jwtSecret)
}
