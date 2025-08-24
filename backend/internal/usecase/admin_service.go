package usecase

import (
	"context"
	"database/sql"
	"errors"
	"url-shortener/backend/internal/auth"
	domain "url-shortener/backend/internal/domain"
)

type AdminService struct {
	userRepo UserRepository
}

func NewAdminService(userRepo UserRepository) *AdminService {
	if userRepo == nil {
		panic("UserRepository cannot be nil")
	}
	return &AdminService{userRepo: userRepo}
}

// Login authenticates a user with email and password, returns JWT token
func (s *AdminService) Login(ctx context.Context, email, password string) (string, int64, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, ErrUnauthorized
		}
		return "", 0, err
	}

	// Verify password
	if user.Password == nil {
		return "", 0, ErrUnauthorized
	}

	if err := auth.VerifyPassword(*user.Password, password); err != nil {
		return "", 0, ErrUnauthorized
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", 0, err
	}

	return token, user.ID, nil
}

// Register creates a new user with hashed password and returns JWT token
func (s *AdminService) Register(ctx context.Context, email, password string) (string, int64, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return "", 0, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return "", 0, err
	}

	// Create user
	user := &domain.User{
		Email:    email,
		Password: &hashedPassword,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return "", 0, err
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", 0, err
	}

	return token, user.ID, nil
}

// ValidateJWT validates a JWT token and returns user ID
func (s *AdminService) ValidateJWT(tokenString string) (int64, error) {
	userID, err := auth.ValidateJWT(tokenString)
	if err != nil {
		return 0, ErrUnauthorized
	}
	return userID, nil
}
