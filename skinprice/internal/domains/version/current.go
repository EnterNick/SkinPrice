package version

import "time"

type Current struct {
	Version    string
	Entrypoint string
	UpdatedAt  time.Time
	Previous   string
}
