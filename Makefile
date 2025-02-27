# Run tests.
test:
	bash -c "set -m; bash '$(CURDIR)/scripts/test.sh'"

# Check code quality.
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run
	npx prettier . --check

format-main:
	go mod tidy
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix
	npx prettier . --write

format-pgbun:
	cd "$(CURDIR)/pgbun" && go mod tidy
	cd "$(CURDIR)/pgbun" && go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix

# Reformat code so it passes the code style lint checks.
format: format-main format-pgbun
