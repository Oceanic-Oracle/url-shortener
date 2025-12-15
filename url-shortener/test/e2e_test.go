package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	app "shortener/internal/bootstrap"
	"shortener/internal/config"
	"shortener/internal/dto"
	"shortener/internal/infra/database"
	"shortener/internal/infra/logger"
)

var cfg = config.MustLoad("./.env")

func TestMain(m *testing.M) {
	logger := logger.SetupLogger(cfg.Env)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := waitForDB(ctx, logger); err != nil {
		log.Fatal("DB not ready:", err)
	}

	go func() {
		app := app.NewBootstrap(cfg, logger)
		app.Run()
	}()

	// TODO: Написать эндпоинт /ping для корректного ожидания
	time.Sleep(10 * time.Second)

	code := m.Run()
	os.Exit(code)
}

func waitForDB(ctx context.Context, log *slog.Logger) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			pool, err := database.GetRedisConnectionPool(ctx, cfg.Storage.URL, log)
			if err != nil {
				continue
			}
			if _, err := pool.Ping(ctx).Result(); err == nil {
				return nil
			}
		}
	}
}

func Request(reqBody []byte, method, path string) ([]byte, int, string, error) {
	url := "http://" + cfg.HTTP.Host + cfg.HTTP.Addr + path
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to create request: %w", err)
	}

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, resp.Header.Get("Location"), fmt.Errorf("failed to read response body: %w", err)
	}

	return body, resp.StatusCode, resp.Header.Get("Location"), nil
}

// =======================================================================================================

func TestCreateCode(t *testing.T) {
	testCases := []struct {
		url          string
		expectStatus int
	}{
		{"https://github.com/Oceanic-Oracle/url-shortener  ", http.StatusCreated},
		{"https://www.youtube.com/watch?v=aS1cJfQ-LrQ&pp=ugUEEgJlbg%3D%3D  ", http.StatusCreated},
		{"https://t.me/Waveeen  ", http.StatusCreated},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			originalURL := strings.TrimSpace(tc.url)

			reqBody := dto.CreateCodeURLRequest{URL: originalURL}
			reqBytes, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatalf("failed to marshal request: %v", err)
			}

			res, statusCode, _, err := Request(reqBytes, "POST", "/shorten")
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if statusCode != tc.expectStatus {
				t.Fatalf("expected status %d, got %d", tc.expectStatus, statusCode)
			}

			var resp dto.CreateCodeURLResponse
			if err := json.Unmarshal(res, &resp); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			if len(resp.Code) != 10 {
				t.Errorf("expected code length 10, got %d", len(resp.Code))
			}

			_, redirectStatus, location, err := Request(nil, "GET", "/"+resp.Code)
			if err != nil {
				t.Fatalf("redirect request failed: %v", err)
			}

			if redirectStatus != http.StatusFound {
				t.Errorf("expected 302 redirect, got %d", redirectStatus)
			}

			if location != originalURL {
				t.Errorf("expected redirect to %q, got %q", originalURL, location)
			}
		})
	}
}

// =======================================================================================================

func TestRedirectNonExistent(t *testing.T) {
	_, statusCode, _, err := Request(nil, "GET", "/nonexistent123")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if statusCode != http.StatusNotFound {
		t.Errorf("expected 404 for non-existent code, got %d", statusCode)
	}
}
