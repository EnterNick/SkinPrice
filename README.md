# SkinPrice

Локальное desktop-приложение на `Wails + Go + React` для ручного отслеживания цен скинов через Steam Market.

## Что умеет v1

- искать скины по названию;
- сохранять их в локальный список;
- обновлять цену одного скина или всего списка;
- удалять скины из списка;
- хранить данные локально в `SQLite` по умолчанию.

## Локальный запуск

1. Скопируйте `.env.example` в `.env` при необходимости.
2. Для desktop runtime по умолчанию достаточно `SQLite`.
3. Установите frontend зависимости:

```bash
cd skinprice/frontend
npm install
```

4. Запустите приложение:

```bash
cd skinprice
wails dev
```

## Переменные окружения

- `APP_DB_DRIVER=sqlite3`
- `APP_DB_NAME=./skinprice.db`
- `STEAM_BASE_URL=https://steamcommunity.com/market`
- `HTTP_TIMEOUT_SECONDS=10`

`Postgres` можно использовать для разработки, но для первого релиза основным storage считается `SQLite`.

## Проверки

```bash
go test ./...
cd skinprice/frontend && npm run lint && npm run build
```
