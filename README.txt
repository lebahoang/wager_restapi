# wager restapi

To start the application:
```
bash ./start.sh
```

To test the connection:
```
curl -X GET http://localhost:8080/ping
```

You should see
```
{
    "message":"pong"
}
```
Perform unit test:
```
go test -v ./...
```

Perform integration test:
```
bash ./start.sh && POSTGRES_URL=postgres://postgres:12345@localhost:55432/postgres?sslmode=disable go test -v -tags="integration" -count=1 ./integration_test
```