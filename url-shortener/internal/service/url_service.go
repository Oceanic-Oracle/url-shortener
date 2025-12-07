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
	shortURL, err := su.repos.URL.GetShortURLByLong(ctx, url.LongURL(longURL))
	if err == nil {
		return string(shortURL), nil
	}

	if !errors.Is(err, url.ErrURLNotFound) {
		return "", WrapErrInternalServer(err)
	}

	for range 5 {
		shortURL := GenerateShortCode(su.cfg.Length)

		err := su.repos.URL.SaveURL(ctx, url.ShortURL(shortURL), url.LongURL(longURL), su.cfg.TTL)
		if err != nil {
			if !errors.Is(err, url.ErrLongURLExists) && !errors.Is(err, url.ErrShortURLExists) {
				return "", WrapErrInternalServer(err)
			} else {
				continue
			}
		}

		return shortURL, nil
	}

	return "", WrapErrInternalServer(ErrFailedGenerateShortCode)
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
