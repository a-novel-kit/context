# Run tests.
test:
	bash -c "set -m; bash '$(CURDIR)/scripts/test.sh'"

format-main:
	npx prettier . --write

format-pgbun:
	cd "$(CURDIR)/pgbun" && go mod tidy
	cd "$(CURDIR)/pgbun" && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run --fix

format-sentry:
	cd "$(CURDIR)/sentry" && go mod tidy
	cd "$(CURDIR)/sentry" && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6 run --fix

# Reformat code so it passes the code style lint checks.
format: format-main format-pgbun format-sentry
