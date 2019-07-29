GOSS_VERSION := 0.3.6

all: build

build: bin
	cd .. && dep ensure
	go build -o bin/demo ..

bin:
	mkdir -p bin
	chmod 777 bin

vendor:
	mkdir -p vendor

clean:
	rm -rf bin vendor

docker-build:
	docker run --rm \
		-v ~/.cache:/.cache \
		-v `pwd`/..:/go/src/github.com/factorysh/alertmanager-multitenant \
		-w /go/src/github.com/factorysh/alertmanager-multitenant \
		-u `id -u` \
		bearstech/golang-dep \
		make -C demo

demo: docker-build
	docker-compose up -d

func_test/goss:
	mkdir -p test_goss/bin
	curl -o test_goss/bin/goss -L https://github.com/aelsabbahy/goss/releases/download/v${GOSS_VERSION}/goss-linux-amd64
	chmod 755 test_goss/bin/goss

func_test: docker-build func_test/goss
	docker-compose down --remove-orphans
	docker-compose up -d alertmanager mailhog demo
	sleep 2
	docker-compose up --abort-on-container-exit --exit-code-from client client
	docker-compose down --remove-orphans
