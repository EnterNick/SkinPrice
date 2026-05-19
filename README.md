# SkinPrice

`SkinPrice` is a desktop application for tracking skin prices with a local-first setup. The app is built on `Wails`, `Go`, and `React`, and stores data in `SQLite` by default.

## Features

- Search skins by name.
- Save skins to a personal watchlist.
- Refresh the price for one skin or the whole list.
- Store data locally without requiring external infrastructure.
- Use `SQLite` by default, with optional `Postgres` support for development scenarios.

## Tech Stack

- Backend: `Go`
- Desktop runtime: `Wails v2`
- Frontend: `React + TypeScript + Vite`
- Default database: `SQLite`

## Release Downloads

Tagged releases publish ready-to-use application archives for all supported desktop platforms.

- Latest release page: [github.com/EnterNick/SkinPrice/releases/latest](https://github.com/EnterNick/SkinPrice/releases/latest)
- Linux (`amd64`): [SkinPrice-linux-amd64.tar.gz](https://github.com/EnterNick/SkinPrice/releases/latest/download/SkinPrice-linux-amd64.tar.gz)
- Windows (`amd64`): [SkinPrice-windows-amd64.zip](https://github.com/EnterNick/SkinPrice/releases/latest/download/SkinPrice-windows-amd64.zip)
- macOS (`universal`): [SkinPrice-macos-universal.zip](https://github.com/EnterNick/SkinPrice/releases/latest/download/SkinPrice-macos-universal.zip)

These links point to assets from the most recent published GitHub Release.

## Project Layout

```text
.
├── skinprice/                # Wails application
│   ├── frontend/             # React + TypeScript frontend
│   ├── internal/             # Backend application code
│   ├── build/                # Wails build configuration and output
│   ├── app.go
│   ├── main.go
│   └── wails.json
├── Makefile
├── go.mod
└── .github/workflows/        # CI/CD pipelines
```

## Requirements

- `Go 1.26.3`
- `Node.js 20`
- `npm`
- `Wails v2`

For Linux desktop builds, system packages required by `Wails` must also be installed.

## Local Development

1. Create the environment file:

```bash
cp .env.example .env
```

2. Install frontend dependencies:

```bash
cd skinprice/frontend
npm install
```

3. Start the desktop app in development mode:

```bash
cd ../
wails dev
```

If `wails` is not installed:

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
```

## Environment Variables

Example configuration is provided in [`.env.example`](.env.example).

Important variables:

- `APP_ENV=local`
- `APP_DB_DRIVER=sqlite3`
- `APP_DB_NAME=./skinprice.db`
- `STEAM_BASE_URL=https://steamcommunity.com/market`
- `LISSKINS_BASE_URL=https://lis-skins.ru`
- `HTTP_TIMEOUT_SECONDS=10`
- `CACHE_TTL_SECONDS=300`
- `LOG_LEVEL=debug`
- `LOG_TO_FILE=true`

`SQLite` is the default and recommended storage for local usage. `Postgres` can be enabled for development if needed.

## Quality Checks

Backend tests:

```bash
go test ./...
```

Frontend checks:

```bash
cd skinprice/frontend
npm run lint
npm run build
```

## Useful Commands

```bash
make deps
make test
make lint-ci
make wails
```

## CI/CD

The main workflow is defined in [`.github/workflows/build.yml`](.github/workflows/build.yml).

- `pull_request`: runs project checks.
- `push` to `master`: runs project checks for the main branch state.
- `push` of tag `v*`: runs checks, builds desktop applications for `Linux`, `Windows`, and `macOS`, then creates a GitHub Release with attached archives.
- `workflow_dispatch`: allows manual desktop builds with downloadable Action artifacts.

Generated release artifact names:

- `SkinPrice-linux-amd64.tar.gz`
- `SkinPrice-windows-amd64.zip`
- `SkinPrice-macos-universal.zip`
