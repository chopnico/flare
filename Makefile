build:
	go mod download
	go build -o build/flare cmd/main.go
clean:
	rm -rf build
	go clean
	rm go.sum
