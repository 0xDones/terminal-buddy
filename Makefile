VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: build
build:
	go build -ldflags "-X main.version=$(VERSION)" -o ./bin/tb .

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: tidy
tidy:
	go mod tidy

install: tidy vet fmt build
	chmod +x ./bin/tb
	mv ./bin/tb /usr/local/bin/tb
