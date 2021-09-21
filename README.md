# squirrel-bot ![develop](https://github.com/mebr0/squirrel-bot/actions/workflows/develop.yml/badge.svg)

Telegram bot for playing card game called `Belka` (eng. Squirrel)

## Variables

Use these variables to run project in `.env` file

```dotenv
LOG_LEVEL=INFO

TELEGRAM_BOT_TOKEN=<token>

DB_HOST=<host>
DB_PORT=<port
DB_NAME=<database>
DB_USER=<username>
DB_PASSWORD=<password>
DB_SSL_MODE=<mode>

GAME_SPEED=<speed>
```

## Commands

`make fmt` - format whole project with gofmt _(do it before any commit)_

`make cover` - run unit tests and show coverage report

`make build` - build project

`make run` - build and run project

## Docker

Use dockerfiles in `build` directory for building images and running containers

Use `build/Dockerfile` for building images on unix systems.
Use `build/Dockerfile.multi` for building images on non-unix systems
