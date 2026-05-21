package httpx

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"syscall"
	"time"
)

type RetryConfig struct {
	Attempts int
	Delay    time.Duration
}

func DoWithRetry(ctx context.Context, client *http.Client, buildRequest func(context.Context) (*http.Request, error), config RetryConfig) (*http.Response, error) {
	attempts := config.Attempts
	if attempts <= 0 {
		attempts = 1
	}
	delay := config.Delay
	if delay <= 0 {
		delay = 500 * time.Millisecond
	}

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		req, err := buildRequest(ctx)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		if !ShouldRetry(err) || attempt == attempts {
			return nil, err
		}

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}

	return nil, lastErr
}

func ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	if errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET) {
		return true
	}

	text := strings.ToLower(err.Error())
	transientMarkers := []string{
		"tls handshake timeout",
		"handshake timeout",
		"connection reset by peer",
		"unexpected eof",
		"server closed idle connection",
		"no such host",
		"temporary failure in name resolution",
	}
	for _, marker := range transientMarkers {
		if strings.Contains(text, marker) {
			return true
		}
	}

	return false
}
