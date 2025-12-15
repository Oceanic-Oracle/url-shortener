package service

import (
	"context"
	"errors"
	"log/slog"

	"shortener/internal/config"
	"shortener/internal/logctx"
	"shortener/internal/repo"
	"shortener/internal/repo/url"
)

type ServiceURL struct {
	cfg   config.URLShortener
	repos *repo.Repo
	log   *slog.Logger
}

func (su *ServiceURL) CreateShortURL(ctx context.Context, longURL string) (string, error) {
	su.log.DebugContext(ctx, "Creating short URL")

	shortURL := GenerateShortCode(longURL, su.cfg.Salt, su.cfg.Length)
	ctx = logctx.WithCode(ctx, shortURL)

	err := su.repos.URL.SaveURL(ctx, url.ShortURL(shortURL), url.LongURL(longURL), su.cfg.TTL)
	if err == nil {
		su.log.DebugContext(ctx, "Short URL saved successfully")
		return shortURL, nil
	}

	if errors.Is(err, url.ErrURLExists) {
		su.log.WarnContext(ctx, "Short code already exists, checking for collision", slog.Any("err", err))

		exLongURL, err := su.repos.URL.GetLongURLByShort(ctx, url.ShortURL(shortURL))
		if err != nil {
			su.log.ErrorContext(ctx, "Failed to fetch existing URL after collision", slog.Any("err", err))
			return "", ErrStorage
		}

		if exLongURL != url.LongURL(longURL) {
			su.log.ErrorContext(ctx, "Short code collision detected", slog.Any("err", err))
			return "", ErrURLCollision
		}

		return shortURL, nil
	}

	su.log.ErrorContext(ctx, "Failed to save short URL", "err", err)

	return "", ErrStorage
}

func (su *ServiceURL) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	su.log.DebugContext(ctx, "Getting URL")

	longURL, err := su.repos.URL.GetLongURLByShortWithTTLUpdate(ctx, url.ShortURL(shortURL), su.cfg.TTL)
	if err != nil {
		if errors.Is(err, url.ErrURLNotFound) {
			su.log.WarnContext(ctx, "URL not found", slog.Any("err", err))
			return "", ErrURLNotFound
		}

		su.log.ErrorContext(ctx, "Failed to get URL", slog.Any("err", err))

		return "", ErrStorage
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
