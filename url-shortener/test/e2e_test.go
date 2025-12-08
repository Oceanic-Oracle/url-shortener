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
	app "shortener/internal/bootstrap"
	"shortener/internal/config"
	"shortener/internal/dto"
	"shortener/internal/infra/database"
	"shortener/internal/infra/logger"
	"testing"
	"time"
)

var (
	env      = "error"
	httpHost = "http://localhost"
	httpAddr = ":9090"
	redisCfg = config.Storage{
		Type: "redis",
		URL:  "redis://:sdnsfnsdnsgqerqew234whdnd@localhost:6380",
	}
)

func TestMain(m *testing.M) {
	cfg := config.MustLoad()
	cfg.Env = env
	cfg.Storage = redisCfg
	cfg.HTTP.Addr = httpAddr

	logger := logger.SetupLogger(cfg.Env)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := waitForDB(ctx, logger); err != nil {
		log.Fatal("DB not ready:", err)
		return
	}

	go func() {
		app := app.NewBootstrap(cfg, logger)
		app.Run()
	}()

	time.Sleep(10 * time.Second)

	code := m.Run()
	os.Exit(code)
}

func waitForDB(ctx context.Context, log *slog.Logger) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			pool, err := database.GetRedisConnectionPool(ctx, redisCfg.URL, log)
			if err != nil {
				time.Sleep(3 * time.Second)
				continue
			}

			_, err = pool.Ping(ctx).Result()
			if err == nil {
				return nil
			}
		}
	}
}

func Request(reqBody []byte, method, path string) ([]byte, int, error) {
	req, err := http.NewRequest(method, httpHost+httpAddr+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, resp.StatusCode, nil
}

// =======================================================================================================

func TestCreateCode(t *testing.T) {
	req := []dto.CreateCodeURLRequest{
		{
			URL: "https://github.com/Oceanic-Oracle/url-shortener",
		},
		{
			URL: "https://www.youtube.com/watch?v=aS1cJfQ-LrQ&pp=ugUEEgJlbg%3D%3D",
		},
		{
			URL: "https://t.me/Waveeen",
		},
	}

	exRes := []dto.CreateCodeURLResponse{
		{
			Code: "I5rVzfYOJx",
		},
		{
			Code: "1lW2Y2l7Xt",
		},
		{
			Code: "G8USrVXmC2",
		},
	}

	status := []int{
		201, 201, 201,
	}

	for i := range req {
		reqBytes, err := json.Marshal(req[i])
		if err != nil {
			t.Errorf("failed to decode request body: %v", err)
			continue
		}

		res, statusCode, err := Request(reqBytes, "POST", "/shorten")
		if err != nil {
			t.Errorf("Request failed: %v", err)
			continue
		}

		if statusCode != status[i] {
			t.Errorf("Incorrect status code. Expected %d, but got %d", status[i], statusCode)
			continue
		}

		var resJson dto.CreateCodeURLResponse
		if err = json.Unmarshal(res, &resJson); err != nil {
			t.Errorf("Failed to unmarshal PR response: %v", err)
			continue
		}

		if resJson != exRes[i] {
			t.Errorf("Create code failed for url. Expected %s, but got %s", exRes[i].Code, resJson.Code)
		}
	}
}

// =======================================================================================================

func TestRedirect(t *testing.T) {
	req := []string{
		"I5rVzfYOJx",
		"1lW2Y2l7Xt",
		"G8USrVXmC2",
	}

	exRes := []int{
		http.StatusFound,
		http.StatusFound,
		http.StatusFound,
	}

	for i := range req {
		_, statusCode, err := Request(nil, "GET", "/"+req[i])
		if err != nil {
			t.Errorf("Request failed: %v", err)
			continue
		}

		if statusCode != exRes[i] {
			t.Errorf("Incorrect status code. Expected %d, but got %d", exRes[i], statusCode)
			continue
		}
	}
}
