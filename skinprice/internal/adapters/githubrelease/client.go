package githubrelease

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	appversion "SkinPrice/skinprice/internal/application/version"
	"SkinPrice/skinprice/internal/shared/httpx"
)

type Client struct {
	BaseURL string
	Repo    string
	HTTP    *http.Client
}

type releaseResponse struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
}

func (c Client) GetLatestRelease(ctx context.Context) (appversion.ReleaseMeta, error) {
	if c.Repo == "" {
		return appversion.ReleaseMeta{}, fmt.Errorf("github repo is required")
	}

	baseURL := strings.TrimRight(c.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}

	resp, err := httpx.DoWithRetry(ctx, c.httpClient(), func(ctx context.Context) (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/repos/"+c.Repo+"/releases/latest", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "skinprice-launcher")
		return req, nil
	}, httpx.RetryConfig{
		Attempts: 3,
		Delay:    750 * time.Millisecond,
	})
	if err != nil {
		return appversion.ReleaseMeta{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return appversion.ReleaseMeta{}, fmt.Errorf("github latest release request failed: %s", resp.Status)
	}

	var decoded releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return appversion.ReleaseMeta{}, err
	}

	meta := appversion.ReleaseMeta{TagName: decoded.TagName}
	for _, asset := range decoded.Assets {
		meta.Assets = append(meta.Assets, appversion.ReleaseAsset{
			Name:        asset.Name,
			DownloadURL: asset.BrowserDownloadURL,
			Size:        asset.Size,
		})
	}

	return meta, nil
}

func (c Client) FindAsset(release appversion.ReleaseMeta, name string) (appversion.ReleaseAsset, error) {
	for _, asset := range release.Assets {
		if asset.Name == name {
			return asset, nil
		}
	}
	return appversion.ReleaseAsset{}, fmt.Errorf("release asset not found: %s", name)
}

func (c Client) httpClient() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	return http.DefaultClient
}
