fmt:
	go fmt ./...

vet:
	go vet ./...

install:
	go get -d ./...

run: fmt vet
	go run ./...

test:
	go test -v ./gateway/...
