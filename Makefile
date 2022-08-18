nebb:
	go build -o bin/nebb cmd/nebb/*.go

fmt:
	go fmt ./...

vet:
	go vet ./...

test: fmt vet
	go test ./... -coverprofile cover.out -short