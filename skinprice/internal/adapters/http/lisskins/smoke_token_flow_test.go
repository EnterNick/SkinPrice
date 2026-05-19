package lisskins

import (
	"SkinPrice/skinprice/internal/adapters/database"
	sourcestate "SkinPrice/skinprice/internal/adapters/database/sourcestate"
	"SkinPrice/skinprice/internal/application"
	appskins "SkinPrice/skinprice/internal/application/skins"
	sharedcrypto "SkinPrice/skinprice/internal/shared/utils/crypto"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLisSkinsTokenSmokeFlow(t *testing.T) {
	conn, err := database.New(&database.Config{Driver: "sqlite3", DBName: t.TempDir() + "/skinprice.db"})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			t.Fatalf("close db: %v", closeErr)
		}
	}()
	if err = database.EnsureSchema(conn); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}

	stateStorage := &sourcestate.Storage{Conn: conn}
	hasUC := appskins.HasLisSkinsToken{Storage: stateStorage}
	hasToken, err := hasUC.Execute()
	if err != nil || hasToken {
		t.Fatalf("new user must have no token: has=%v err=%v", hasToken, err)
	}

	key := base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef"))
	cipher, err := sharedcrypto.NewTokenCipherFromBase64(key)
	if err != nil {
		t.Fatalf("cipher init: %v", err)
	}
	saveUC := appskins.SaveLisSkinsToken{Storage: stateStorage, Cipher: cipher}
	getUC := appskins.GetLisSkinsToken{Storage: stateStorage, Cipher: cipher}

	const plainToken = "tok_test_123456789"
	if err = saveUC.Execute(plainToken); err != nil {
		t.Fatalf("save token: %v", err)
	}
	rawEncrypted, err := stateStorage.GetLisSkinsToken()
	if err != nil {
		t.Fatalf("read encrypted: %v", err)
	}
	if rawEncrypted == plainToken || strings.Contains(rawEncrypted, plainToken) {
		t.Fatalf("token must be encrypted in DB")
	}

	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if got := r.Header.Get(lisSkinsAuthHeader); got != buildLisSkinsAuthHeaderValue(plainToken) {
			t.Fatalf("bad auth header: %q", got)
		}
		if calls == 1 {
			_, _ = w.Write([]byte(`{"data":[],"total_count":0}`))
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	httpStorage := &Storage{Client: server.Client(), BaseURL: server.URL, TokenProvider: getUC}
	pagination := &application.Pagination{Limit: 20, Offset: 0}
	if _, err = httpStorage.GetList(appskins.SearchCriteria{}, pagination); err != nil {
		t.Fatalf("search should succeed: %v", err)
	}
	_, err = httpStorage.GetList(appskins.SearchCriteria{}, pagination)
	if !errors.Is(err, appskins.ErrLisSkinsTokenInvalid) {
		t.Fatalf("expected invalid token branch, got: %v", err)
	}
}
