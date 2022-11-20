.PHONY: build
build:
	go build -o bin/tg2fedi ./cmd/tg2fedi

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/tg2fedi.linux.amd64 ./cmd/tg2fedi

.PHONY: build-docker
build-docker:
	docker build --platform linux/amd64 -t tg2fedi .
