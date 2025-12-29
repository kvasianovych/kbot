APP=$(shell basename $(shell git remote get-url origin))
REGISTRY=kvasianovych
VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS=linux
TARGETARCH=amd64

format:
	gofmt -s -w ./

# lint:
# 	staticcheck ./...

# tools:
# 	go install honnef.co/go/tools/cmd/staticcheck@latest

test:
	go test

get:
	go get

build: format get
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o kbot -ldflags "-X=github.com/kvasianovych/kbot/cmd.appVersion=v${VERSION}"

image:
	docker build -t ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH} .

push:
	docker push ${REGISTRY}/${APP}:${VERSION}-${TARGETARCH}

clean:
	rm kbot
