package usecase

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"os"
	"strconv"
	domain "url-shortener/backend/internal/domain"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const codeLength = 6
const maxGenShortCodeRetries = 5
const maxGenAPIKeyRetries = 5
const apiKeyLength = 32

var (
	ErrLinkAlreadyExists           = errors.New("link already exists for this user")
	ErrMaxShortCodeRetriesExceeded = errors.New("could not generate a unique short code after multiple retries")
	ErrMaxAPIKeyRetriesExceeded    = errors.New("could not generate a unique API key after multiple retries")
	ErrLinkNotFound                = errors.New("short code not found")
	ErrLinkLimitExceeded           = errors.New("free plan link limit reached")
	ErrShortCodeAlreadyExists      = errors.New("custom short code already exists")
	ErrUnauthorized                = errors.New("unauthorized access")
	ErrUserNotFound                = errors.New("user not found")
	ErrUserHasNoAPIKey             = errors.New("user has no API key")
)

// FreePlanMaxLinks is configurable via env FREE_PLAN_MAX_LINKS (default 10)
var FreePlanMaxLinks = func() int {
	if v := os.Getenv("FREE_PLAN_MAX_LINKS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 10
}()

type LinkRepository = domain.LinkRepository
type UserRepository = domain.UserRepository

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

func (s *ShortenerService) CreateShortLink(ctx context.Context, apiKeyID int64, longURL string, customShortCode string, password string) (*domain.Link, error) {

	var shortCode string
	if customShortCode != "" {
		// Check if custom shortCode already exists
		existingLink, err := s.linkRepo.GetLinkByShortCode(ctx, customShortCode)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if existingLink != nil {
			return nil, ErrShortCodeAlreadyExists
		}
		shortCode = customShortCode
	} else {
		// Generate random shortCode
		for i := 0; i < maxGenShortCodeRetries; i++ {
			tmpCode := generateRandomCode(codeLength)

			link, err := s.linkRepo.GetLinkByShortCode(ctx, tmpCode)
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}

			if link == nil {
				shortCode = tmpCode
				break
			}
		}
		if shortCode == "" {
			return nil, ErrMaxShortCodeRetriesExceeded
		}
	}

	var apiKeyIDPtr *int64
	var passwordPtr *string
	if apiKeyID != 0 {
		apiKeyIDPtr = &apiKeyID
	}
	if password != "" {
		passwordPtr = &password
	}
	link := &domain.Link{
		APIKeyID:  apiKeyIDPtr,
		LongURL:   longURL,
		ShortCode: shortCode,
		Password:  passwordPtr,
	}

	if err := s.linkRepo.CreateLink(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *ShortenerService) ResolveLink(ctx context.Context, shortCode string) (string, error) {
	link, err := s.linkRepo.GetLinkByShortCode(ctx, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrLinkNotFound
		}
		return "", err
	}

	if err := s.linkRepo.TrackClick(ctx, shortCode); err != nil {
		return "", err
	}

	return link.LongURL, nil
}

func (s *ShortenerService) CreateUser(ctx context.Context, email string, password *string) (*domain.User, error) {
	user := &domain.User{
		Email:    email,
		Password: password,
	}
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *ShortenerService) CreateAPIKey(ctx context.Context, userID int64) (string, error) {
	key := s.generateUniqueAPIKey(ctx)
	if key == "" {
		return "", ErrMaxAPIKeyRetriesExceeded
	}

	if err := s.userRepo.CreateAPIKey(ctx, userID, key); err != nil {
		return "", err
	}
	return key, nil
}
func generateRandomCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func (s *ShortenerService) generateUniqueAPIKey(ctx context.Context) string {
	for i := 0; i < maxGenAPIKeyRetries; i++ {
		apiKey := generateRandomCode(apiKeyLength)
		_, err := s.userRepo.GetAPIKeyIDByAPIKey(ctx, apiKey)
		if err == sql.ErrNoRows {
			return apiKey
		}
	}
	return ""
}

func (s *ShortenerService) GetShortURLsByAPIKey(ctx context.Context, apiKeyID int64) ([]*domain.Link, error) {
	links, err := s.linkRepo.GetListLinksByAPIKeyID(ctx, apiKeyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}

	return links, nil
}

// GetLinksWithoutAPIKey returns all links that were created without an API key (i.e., APIKeyID is NULL).
func (s *ShortenerService) GetLinkWithoutAPIKey(ctx context.Context, shortCode string) ([]*domain.Link, error) {
	links, err := s.linkRepo.GetLinkNoAPIKey(ctx, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrLinkNotFound
		}
		return nil, err
	}
	return links, nil
}

func (s *ShortenerService) GetFirstAPIKey(ctx context.Context, userID int64) (string, error) {
	apiKeys, err := s.userRepo.GetAPIKeysByUserID(ctx, userID, 1, 0)
	if err != nil {
		return "", err
	}
	if len(apiKeys) == 0 {
		return "", ErrUserHasNoAPIKey
	}
	return apiKeys[0].Key, nil
}

func (s *ShortenerService) DeleteShortLink(ctx context.Context, apiKey string, shortCode string) error {
	apiKeyID, err := s.userRepo.GetAPIKeyIDByAPIKey(ctx, apiKey)
	if err != nil {
		return ErrUnauthorized
	}

	err = s.linkRepo.SoftDeleteByShortCode(ctx, apiKeyID, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrLinkNotFound
		}
		return err
	}
	return nil
}

func (s *ShortenerService) GetLinksByUser(ctx context.Context, userID int64) ([]*domain.Link, map[int64]string, error) {
	apiKeys, err := s.userRepo.GetAPIKeysByUserID(ctx, userID, 0, 0)
	if err != nil {
		if err == sql.ErrNoRows {
			// User exists but has no API keys, return empty slice
			return []*domain.Link{}, nil, err
		}
		return nil, nil, err
	}

	links := make([]*domain.Link, 0)
	apiKeysMap := make(map[int64]string)
	for _, apiKey := range apiKeys {
		linksByAPIKey, err := s.GetShortURLsByAPIKey(ctx, apiKey.ID)
		if err != nil {
			if err == ErrLinkNotFound {
				continue
			}
			return nil, nil, err
		}
		links = append(links, linksByAPIKey...)
		apiKeysMap[apiKey.ID] = apiKey.Key
	}
	return links, apiKeysMap, nil
}
