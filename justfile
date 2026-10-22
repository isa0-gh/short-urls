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

fmt:
	go fmt ./...

vet:
	go vet ./...

gen-docs: require-swag
	swag init -g cmd/api/main.go --parseInternal -o docs

check: require-govulncheck test vet sqlc-check
    govulncheck ./...

default:
	@just --list

sqlc:
	@just require-sqlc
	sqlc generate

sqlc-check: require-sqlc
	@diff=$(sqlc diff 2>&1); diff_status=$?; \
	if [ -n "$diff" ]; then \
		printf '%s\n' "$diff"; \
		echo "Generated SQL code is out of date. Run 'just sqlc'." >&2; \
		exit 1; \
	fi; \
	if [ $diff_status -ne 0 ]; then \
		echo "sqlc diff failed with exit code $diff_status." >&2; \
		exit $diff_status; \
	fi

tools:
	go install golang.org/x/vuln/cmd/govulncheck@{{GOVULNCHECK_VERSION}}
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@{{SQLC_VERSION}}
	go install github.com/swaggo/swag/cmd/swag@{{SWAG_VERSION}}
	@echo "Ensure $(go env GOPATH)/bin is in your PATH."

require-govulncheck:
	@command -v govulncheck >/dev/null 2>&1 || { echo "govulncheck is not installed; run 'just tools'." >&2; exit 1; }

require-sqlc:
	@command -v sqlc >/dev/null 2>&1 || { echo "sqlc is not installed; run 'just tools'." >&2; exit 1; }

require-swag:
	@command -v swag >/dev/null 2>&1 || { echo "swag is not installed; run 'just tools'." >&2; exit 1; }

build: sqlc-check gen-docs
	go build -o {{BINARY_NAME}} {{MAIN_PATH}}

run: sqlc-check
	go run {{MAIN_PATH}}

release version: sqlc-check
	@current_branch=$(git branch --show-current); \
	if [ "$current_branch" != "main" ]; then \
		echo "Error: You must be on the 'main' branch to tag a release (currently on '$current_branch')"; \
		exit 1; \
	fi
	git tag {{version}}
	git push origin {{version}}

docker-build version: sqlc-check
	docker build -t {{DOCKER_IMAGE}}:{{version}} .

clean:
	rm -f {{BINARY_NAME}}
