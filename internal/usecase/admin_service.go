package usecase

import (
	"context"
	"time"
	"url-shortener/internal/domain"
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

func (s *AdminService) CreateUser(ctx context.Context, email, plan, role string, planExpiresAt *time.Time) (*domain.User, error) {
	user := &domain.User{Email: email, Plan: plan, Role: role, PlanExpiresAt: planExpiresAt}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AdminService) CreateAPIKeyForUser(ctx context.Context, userID int64, key string) error {
	return s.userRepo.CreateAPIKey(ctx, userID, key)
}

func (s *AdminService) SoftDeleteUser(ctx context.Context, userID int64) error {
	return s.userRepo.SoftDeleteByID(ctx, userID)
}

func (s *AdminService) UpdateUserPlan(ctx context.Context, userID int64, plan string, planExpiresAt *time.Time) error {
	return s.userRepo.UpdatePlanAndExpiry(ctx, userID, plan, planExpiresAt)
}
