# migrate

`migrate` is a zero-dependency utility to run SQL migration against PostgreSQL.

Built on top of [goose](https://github.com/pressly/goose).

## download

https://github.com/kofile/migrate/releases

## build

Install dependencies with [`dep`](https://golang.github.io/dep/).

```
dep ensure
```

### distros

**NOTE**: Requires Docker & `tar`.

```sh
make
make archive
```

## author

[neezer](https://github.com/neezer/)
