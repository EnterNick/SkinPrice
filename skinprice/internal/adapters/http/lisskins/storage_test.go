package lisskins

import (
	"SkinPrice/skinprice/internal/application"
	"SkinPrice/skinprice/internal/application/skins"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

type tokenProviderStub struct {
	token string
	err   error
}

func (s tokenProviderStub) Execute() (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.token, nil
}

func TestGetListAddsAuthorizationHeader(t *testing.T) {
	token := "secret-token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := buildLisSkinsAuthHeaderValue(token)
		if got := r.Header.Get(lisSkinsAuthHeader); got != expected {
			t.Fatalf("unexpected auth header: %q", got)
		}
		_, _ = w.Write([]byte(`{"data":[],"total_count":0}`))
	}))
	defer server.Close()

	storage := &Storage{Client: server.Client(), BaseURL: server.URL, TokenProvider: tokenProviderStub{token: token}}
	_, err := storage.GetList(skins.SearchCriteria{}, &application.Pagination{Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("GetList() error = %v", err)
	}
}

func TestGetListWithMissingTokenDoesNotSendRequest(t *testing.T) {
	var requests int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	storage := &Storage{Client: server.Client(), BaseURL: server.URL, TokenProvider: tokenProviderStub{token: ""}}
	_, err := storage.GetList(skins.SearchCriteria{}, &application.Pagination{Limit: 10, Offset: 0})
	if !errors.Is(err, skins.ErrLisSkinsTokenMissing) {
		t.Fatalf("expected ErrLisSkinsTokenMissing, got %v", err)
	}
	if atomic.LoadInt32(&requests) != 0 {
		t.Fatalf("expected no outbound request, got %d", requests)
	}
}

func TestGetListMapsAuthStatusToInvalidToken(t *testing.T) {
	tests := []int{http.StatusUnauthorized, http.StatusForbidden}
	for _, status := range tests {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			}))
			defer server.Close()

			storage := &Storage{Client: server.Client(), BaseURL: server.URL, TokenProvider: tokenProviderStub{token: "token"}}
			_, err := storage.GetList(skins.SearchCriteria{}, &application.Pagination{Limit: 10, Offset: 0})
			if !errors.Is(err, skins.ErrLisSkinsTokenInvalid) {
				t.Fatalf("expected ErrLisSkinsTokenInvalid, got %v", err)
			}
		})
	}
}
