package usecase

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	domain "url-shortener/internal/domain"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const codeLength = 6
const maxRetries = 5

var (
	ErrLinkAlreadyExists  = errors.New("link already exists for this user")
	ErrMaxRetriesExceeded = errors.New("could not generate a unique short code after multiple retries")
	ErrLinkNotFound       = errors.New("short code not found")
)

type LinkRepository interface {
	Create(ctx context.Context, link *domain.Link) error
	FindByShortCode(ctx context.Context, shortCode string) (*domain.Link, error)
	ListByUser(ctx context.Context, userID int64) ([]*domain.Link, error)
	SoftDeleteByShortCode(ctx context.Context, userID int64, shortCode string) error
	TrackClick(ctx context.Context, shortCode string) error

	FindLinkCountByUserIDAndLongURL(ctx context.Context, userID int64, longURL string) (int, error)
}

type UserRepository interface {
	FindByAPIKey(ctx context.Context, apiKey string) (*domain.User, error)
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

	count, err := s.linkRepo.FindLinkCountByUserIDAndLongURL(ctx, userID, longURL)

	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrLinkAlreadyExists
	}
	var shortCode string
	for i := 0; i < maxRetries; i++ {
		tmpCode := generateRandomCode(codeLength)

		link, err := s.linkRepo.FindByShortCode(ctx, tmpCode)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		if link == nil {
			shortCode = tmpCode
			break
		}
	}

	link := &domain.Link{
		UserID:    userID,
		LongURL:   longURL,
		ShortCode: shortCode,
	}
	if err := s.linkRepo.Create(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *ShortenerService) ResolveLink(ctx context.Context, shortCode string) (string, error) {
	link, err := s.linkRepo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	if link == nil {
		return "", ErrLinkNotFound
	}

	if err := s.linkRepo.TrackClick(ctx, shortCode); err != nil {
		return "", err
	}

	return link.LongURL, nil
}

func (s *ShortenerService) ListLinksByUser(ctx context.Context, userID int64) ([]*domain.Link, error) {
	return s.linkRepo.ListByUser(ctx, userID)
}

func (s *ShortenerService) SoftDeleteByCode(ctx context.Context, userID int64, shortCode string) error {
	return s.linkRepo.SoftDeleteByShortCode(ctx, userID, shortCode)
}

func generateRandomCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}
