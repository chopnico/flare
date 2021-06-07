build:
	go mod download
	go build -o build/flare_linux_x64 cmd/main.go
clean:
	rm -rf build
	go clean
	rm go.sum
