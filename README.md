# fyler

![Actions Status](https://github.com/fylerx/fyler/actions/workflows/go.yml/badge.svg)
[![codecov](https://codecov.io/gh/fylerx/fyler/branch/main/graph/badge.svg)](https://codecov.io/gh/fylerx/fyler)

Run Fylerx Dev Environment

```
docker-compose up
```

Run Tests

```
go test -v ./...
```

Create migration

```
migrate create -ext sql -dir db/migrations -seq create_users_table
```

Installing lefthook

```
brew install lefthook
```

```
MY_FAKTORY_URL=tcp://:qwerty@localhost:7419 FAKTORY_PROVIDER=MY_FAKTORY_URL go run cmd/dispatcher/main.go --race


MY_FAKTORY_URL=tcp://:qwerty@localhost:7419 FAKTORY_PROVIDER=MY_FAKTORY_URL go run cmd/worker/main.go --race



go run cmd/cli/main.go new test_project https://webhook.site/7c00b9e9-c7a2-408a-bbd1-9e5215086e76
```
