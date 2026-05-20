package filedownloader

import (
	"context"
	"fmt"
	"io"
	"net/http"

	appversion "SkinPrice/skinprice/internal/application/version"
)

type Downloader struct {
	HTTP *http.Client
}

func (d Downloader) Download(ctx context.Context, asset appversion.ReleaseAsset) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, asset.DownloadURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "skinprice-launcher")

	resp, err := d.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("asset download failed: %s", resp.Status)
	}
	return resp.Body, nil
}

func (d Downloader) httpClient() *http.Client {
	if d.HTTP != nil {
		return d.HTTP
	}
	return http.DefaultClient
}
