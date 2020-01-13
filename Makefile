export BIN_DIR=bin
export ADDRS=http://0.0.0.0:9000
export PORT=8000
export STALE_TIMEOUT=1

export IMAGE_NAME=freundallein/go-lb:latest

init:
	git config core.hooksPath .githooks
run:
	go run main.go
test:
	go test -cover ./...
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o $$BIN_DIR/go-lb
build-healthchecker:
	cd healthchecker && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o ../$$BIN_DIR/healthchecker
dockerbuild:
	docker build -t $$IMAGE_NAME -f Dockerfile .
distribute:
	echo "$$DOCKER_PASSWORD" | docker login -u "$$DOCKER_USERNAME" --password-stdin
	docker build -t $$IMAGE_NAME .
	docker push $$IMAGE_NAME