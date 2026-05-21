package httpx

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestShouldRetryRecognizesTransientErrors(t *testing.T) {
	cases := []error{
		context.DeadlineExceeded,
		errors.New("net/http: TLS handshake timeout"),
		io.EOF,
		timeoutError{},
	}

	for _, err := range cases {
		if !ShouldRetry(err) {
			t.Fatalf("ShouldRetry(%q) = false, want true", err)
		}
	}
}

func TestShouldRetryIgnoresContextCanceled(t *testing.T) {
	if ShouldRetry(context.Canceled) {
		t.Fatalf("ShouldRetry(context.Canceled) = true, want false")
	}
}

func TestDoWithRetryRetriesTransientError(t *testing.T) {
	attempts := 0
	client := &http.Client{
		Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
			attempts++
			if attempts == 1 {
				return nil, errors.New("net/http: TLS handshake timeout")
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("ok")),
			}, nil
		}),
	}

	resp, err := DoWithRetry(context.Background(), client, func(ctx context.Context) (*http.Request, error) {
		return http.NewRequestWithContext(ctx, http.MethodGet, "https://example.com", nil)
	}, RetryConfig{Attempts: 2, Delay: time.Millisecond})
	if err != nil {
		t.Fatalf("DoWithRetry() error = %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

var _ net.Error = timeoutError{}
