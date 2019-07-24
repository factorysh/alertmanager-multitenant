all: build

build: bin
	dep ensure
	go build -o bin/demo

bin:
	mkdir -p bin
	chmod 777 bin

vendor:
	mkdir -p vendor

clean:
	rm -rf bin vendor

docker-build:
	docker run --rm \
		-v `pwd`:/go/src/github.com/factorysh/alertmanager-multitenant \
		-w /go/src/github.com/factorysh/alertmanager-multitenant \
		-u `id -u` \
		bearstech/golang-dep \
		make
