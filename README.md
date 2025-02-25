# jobsity-chat-assignment

A chat application, it uses `WebSocket`s instead of _RabbitMQ_, but, uses the same principle

* SQLite is used as persistence layer
* User registration is implemented
* User log in is implemented
* Cookie storage is used to know if the user is logged in or not
* Multiple users can chat at the same time through the browser
* `chat_service_test.go` file tests concurrently code implemented 

# Help

```sh
$ go run cmd/chat/main.go -help
```

# Run

```sh
$ go run cmd/chat/main.go -debug -sessionKey=supersecret
```

# Test

```sh
$ go test -race -count=1 ./...
```