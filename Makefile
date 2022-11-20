DEFAULT_OUT=bin/tg2fedi
ifdef GOOS
DEFAULT_OUT:=$(DEFAULT_OUT).$(GOOS)
endif
ifdef GOARCH
DEFAULT_OUT:=$(DEFAULT_OUT).$(GOARCH)
endif
ifeq ($(GOOS),windows)
DEFAULT_OUT:=$(DEFAULT_OUT).exe
endif
OUT?=$(DEFAULT_OUT)

.PHONY: build
build:
	$(info $(OUTT))
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(OUT) ./cmd/tg2fedi

.PHONY: build-linux
build-linux:
	$(MAKE) build GOOS=linux GOARCH=amd64 OUT=bin/cmd

.PHONY: lint
lint:
	golangci-lint run



.PHONY: build-docker
build-docker:
	docker build --platform linux/amd64 -t tg2fedi .
