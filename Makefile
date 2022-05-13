install:
	go get -d ./...

run:
	go run src/main.go

test:
	go test -v ./src/gateway/...
