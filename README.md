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


## LisSkins токен

Для поиска через LisSkins нужен персональный API-токен.

### Обязательные env

- `TOKEN_ENCRYPTION_KEY` — base64-ключ длиной 32 байта (AES-256-GCM).

Пример генерации:

```bash
openssl rand -base64 32
```

### Где взять токен LisSkins

1. Откройте профиль LisSkins: `https://lis-skins.com/profile/settings`.
2. Сгенерируйте/скопируйте API-токен в личном кабинете.
3. В приложении переключитесь на источник **LisSkins**, вставьте токен и нажмите **Сохранить токен**.

### Как сбросить токен локально

- Через UI: повторно откройте экран поиска LisSkins и очистите/обновите токен.
- Через БД SQLite вручную:

```bash
sqlite3 skinprice/skinprice.db "DELETE FROM source_states WHERE source = 'lisskins';"
```

После удаления приложение снова покажет CTA на ввод токена.
