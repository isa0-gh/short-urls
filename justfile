set export

GOEXPERIMENT := "jsonv2"
BINARY_NAME := "api"
MAIN_PATH := "cmd/api/main.go"
DOCKER_IMAGE := "github.com/isa0-gh/short-urls"
SQLC_VERSION := "v1.31.1"
SWAG_VERSION := "v1.16.6"
GOVULNCHECK_VERSION := "v1.9.0"

test:
	go test -v ./...

gen-docs:
	swag init -g cmd/api/main.go --parseInternal -o docs

check: test
    govulncheck ./...

default:
	@just --list

sqlc:
	sqlc generate

build: gen-docs
	go build -o {{BINARY_NAME}} {{MAIN_PATH}}

run: sqlc
	go run {{BINARY_NAME}}

release version:
    @current_branch=$$(git branch --show-current); \
    if [ "$$current_branch" != "main" ]; then \
        echo "Error: You must be on the 'main' branch to tag a release (currently on '$$current_branch')"; \
        exit 1; \
    fi
    git tag {{version}}
    git push origin {{version}}

docker-build version:
	docker build -t {{DOCKER_IMAGE}}:{{version}} .

clean:
	rm $BINARY_NAME
