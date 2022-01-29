# fyler
![fyler](https://github.com/fylerx/fyler/actions/workflows/go.yml/badge.svg)

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
