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
