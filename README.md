# Shopping cart API

Implementation of an API Service for an online shop

## Spec

Shopping cart easily add/remove items and handle current promotions

## Assumptions

1.

## Running

Runs with basic config in [cmd/api.shopping/config.go](cmd/api.shopping/config.go).
Optionally, a JSON config file can be passed via the `-config` flag

```sh
go build ./cmd/api.shopping/ && ./api.shopping -config config.json
```

## Unit tests

```sh
go test
```

Or in a Docker container

```sh
make test
```
