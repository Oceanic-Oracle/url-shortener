package service

import (
	"context"
	"errors"
	"log/slog"

	"shortener/internal/config"
	"shortener/internal/repo"
	"shortener/internal/repo/url"
)

type ServiceURL struct {
	cfg   config.URLShortener
	repos *repo.Repo
	log   *slog.Logger
}

func (su *ServiceURL) CreateShortURL(ctx context.Context, longURL string) (string, error) {
	shortURL := GenerateShortCode(longURL, su.cfg.Salt, su.cfg.Length)

	err := su.repos.URL.SaveURL(ctx, url.ShortURL(shortURL), url.LongURL(longURL), su.cfg.TTL)
	if err == nil {
		return shortURL, nil
	}

	if errors.Is(err, url.ErrURLExists) {
		exLongURL, err := su.repos.URL.GetLongURLByShort(ctx, url.ShortURL(shortURL))
		if err != nil {
			return "", WrapErrInternalServer(err)
		}

		if exLongURL != url.LongURL(longURL) {
			return "", WrapErrInternalServer(ErrURLCollision)
		}

		return shortURL, nil
	}

	return "", WrapErrInternalServer(err)
}

func (su *ServiceURL) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	longURL, err := su.repos.URL.GetLongURLByShortWithTTLUpdate(ctx, url.ShortURL(shortURL), su.cfg.TTL)
	if err != nil {
		if errors.Is(err, url.ErrURLNotFound) {
			return "", WrapErrNotFound(err)
		}

		return "", WrapErrInternalServer(err)
	}

	return string(longURL), nil
}

func NewServiceURL(cfg config.URLShortener, repos *repo.Repo, log *slog.Logger) *ServiceURL {
	return &ServiceURL{
		cfg:   cfg,
		repos: repos,
		log:   log,
	}
}
