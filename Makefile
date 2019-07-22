build:
	mkdir -p bin
	dep ensure
	go build -o bin/am_middleware
