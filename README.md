# migrate

[![Build Status](https://ci.neezer.info/api/badges/kofile/migrate/status.svg)](https://ci.neezer.info/kofile/migrate)

`migrate` is a zero-dependency utility to run SQL migration against PostgreSQL.

Built on top of [goose](https://github.com/pressly/goose).

```
$ ./migrate -h
Usage: migrate [OPTIONS] COMMAND
Drivers:
    postgres
Examples:
    migrate status
Available ENV vars:
    DB_HOST                           (required)
    DB_PORT                           (required)
    DB_DATABASE                       (required)
    DB_USERNAME
    DB_PASSWORD
    DB_OPTIONS    eg. sslmode=disable

  -dir string
        directory with migration files (default ".")
  -useEnv bool
        use local .env file

Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version
```

## download

https://github.com/kofile/migrate/releases

## build

Install dependencies with [`dep`](https://golang.github.io/dep/).

```
dep ensure
```

### alpine

1. Copy/mount the source into a `golang` Docker container.
2. Install dependencies (see imports in `main.go`).
3. `CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo .`

### macOS, linux, etc.

Use [gox](https://github.com/mitchellh/gox).

## author

[neezer](https://github.com/neezer/)
