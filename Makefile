HELPERS := ./scripts
OUTPATH := ./bin

EXEC_NAME := pocket

OS := $(shell uname)
GO_SRC := main.go
GOVARS := GO111MODULE=on

ifeq ($(OS),Linux)
	# This is required to get a statically linked binary.
	# Doing this on MacOS breaks something with networking.
	GOVARS := $(GOVARS) CGO_ENABLED=0 GOOS=linux GOARCH=amd64
endif

.PHONY: build build-linux check-clean ci clean dev format install-deps lint start

install-deps:
	go mod download
	go mod tidy

build:
	$(GOVARS) go build \
	-o ${EXEC_NAME} \
	${GO_SRC}
	
build-linux:
	$(GOVARS) go build \
		-ldflags '-extldflags "-static"' \
		-o $(OUTPATH)/${EXEC_NAME} 

ci: install-deps build

start: ci
	./pocket serve --debug --http="$(shell ipconfig getifaddr en0):8090"

check-clean:
	# Ensures working dir is clean
	${HELPERS}/check-clean

format:
	${HELPERS}/format

lint: bin/golangci-lint
	${HELPERS}/lint

clean:
	# Remove files and directories ignored by .gitignore files
	git clean -fdX artifacts bin

docs:
	npm list -g @mermaid-js/mermaid-cli >/dev/null || npm install -g @mermaid-js/mermaid-cli
	mmdc -i docs/diagram.md -o docs/diagram.svg

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.52.2

dev: build
	azurite-blob --blobHost 127.0.0.1 -l artifacts & 
	func start
