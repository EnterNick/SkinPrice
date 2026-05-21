package settings

import "context"

type Storage interface {
	GetAppSettings(ctx context.Context) (AppSettings, error)
	SaveAppSettings(ctx context.Context, settings AppSettings) error
}
