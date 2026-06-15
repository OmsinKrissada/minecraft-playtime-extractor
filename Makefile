.PHONY: all build-web build-go

all: build-web build-go

build-web:
	cd website && pnpm build

build-go:
	go build -ldflags="-s -w" -trimpath

clean:
	go clean
	cd website && rm -rf node_modules dist