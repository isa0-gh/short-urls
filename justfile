set export

GOEXPERIMENT := "jsonv2"
BINARY_NAME := "api"
MAIN_PATH := "cmd/api/main.go"
DOCKER_IMAGE := "short-urls"

test:
	go test -v ./...

check: test
    govulncheck ./...

default:
	@just --list

sqlc:
	sqlc generate

build: sqlc
	go build -o {{BINARY_NAME}} {{MAIN_PATH}}

run: sqlc
	go run {{BINARY_NAME}}
gen-docs:
	swag init -g cmd/api/main.go --parseInternal -o docs
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
