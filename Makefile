all: build

build: bin
	dep ensure
	go build -o bin/am_middleware

bin:
	mkdir -p bin
	chmod 777 bin

vendor:
	mkdir -p vendor

clean:
	rm -rf bin vendor

test: vendor
	go test -v -cover ./...
