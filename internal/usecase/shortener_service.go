package usecase

import (
	"context"
	"math/rand"
	"time"
	domain "url-shortener/internal/domain"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const codeLength = 6
const maxRetries = 5

type LinkRepository interface {
	Create(ctx context.Context, link *domain.Link) error
	GetByShortCode(ctx context.Context, shortCode string) (*domain.Link, error)
	ListByUser(ctx context.Context, userID int64) ([]*domain.Link, error)
	SoftDeleteByShortCode(ctx context.Context, userID int64, shortCode string, deletedAt time.Time) error
	TrackClick(ctx context.Context, shortCode string, at time.Time) error
}

type UserRepository interface {
	GetByAPIKey(ctx context.Context, apiKey string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}

type ShortenerService struct {
	linkRepo LinkRepository
	userRepo UserRepository
}

func NewShortenerService(linkRepo LinkRepository, userRepo UserRepository) *ShortenerService {
	if linkRepo == nil {
		panic("LinkRepository cannot be nil")
	}
	if userRepo == nil {
		panic("UserRepository cannot be nil")
	}
	return &ShortenerService{linkRepo: linkRepo, userRepo: userRepo}
}

func (s *ShortenerService) CreateShortLink(ctx context.Context, userID int64, longURL string) (*domain.Link, error) {
	//Currently I'm mocking
	shortCode := generateRandomCode(codeLength)
	// exo

	link := &domain.Link{
		UserID:    userID,
		LongURL:   longURL,
		ShortCode: shortCode,
		CreatedAt: time.Now(),
	}
	if err := s.linkRepo.Create(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *ShortenerService) ResolveLink(ctx context.Context, shortCode string) (string, error) {
	link, err := s.linkRepo.GetByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	now := time.Now()
	if err := s.linkRepo.TrackClick(ctx, shortCode, now); err != nil {
		return "", err
	}

	return link.LongURL, nil
}

func (s *ShortenerService) ListLinksByUser(ctx context.Context, userID int64) ([]*domain.Link, error) {
	return s.linkRepo.ListByUser(ctx, userID)
}

func (s *ShortenerService) SoftDeleteByCode(ctx context.Context, userID int64, shortCode string) error {
	return s.linkRepo.SoftDeleteByShortCode(ctx, userID, shortCode, time.Now())
}

func generateRandomCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}
