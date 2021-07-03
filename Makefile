build:
	go mod vendor
	go fmt ./...
	go build -o build/flare/flare cmd/flare/flare.go
clean:
	rm -rf build
	go clean
fmt:
	go fmt ./...
