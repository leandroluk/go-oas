GO          ?= go
MODULE      := github.com/leandroluk/go-oas
PKGS        := ./v2 ./v3 ./v3_1
COVERFILE   := coverage.txt
COVERHTML   := coverage.html

.PHONY: deps build test test.ci test.html test.func coverage.save fmt vet lint clean help \
        test.v2 test.v3 test.v3_1 coverage.v2 coverage.v3 coverage.v3_1

## Download deps
deps:
	$(GO) mod download

## Build all versions
build:
	$(GO) build ./...

## Run all tests
test:
	$(GO) test ./...

## Run tests with coverage (global)
test.ci:
	$(GO) test -coverpkg=./... -coverprofile="$(COVERFILE)" -covermode=atomic ./...

## Open coverage in browser
test.html: test.ci
	$(GO) tool cover -html="$(COVERFILE)"

## Show coverage by function
test.func: test.ci
	$(GO) tool cover -func="$(COVERFILE)"

## Save coverage report
coverage.save: test.ci
	$(GO) tool cover -html="$(COVERFILE)" -o "$(COVERHTML)"
	@echo "Report saved at: $(COVERHTML)"

## -------------------------
## Per-version helpers
## -------------------------

test.v2:
	$(GO) test ./v2/...

test.v3:
	$(GO) test ./v3/...

test.v3_1:
	$(GO) test ./v3_1/...

coverage.v2:
	$(GO) test -coverpkg=./v2 -coverprofile="$(COVERFILE)" -covermode=atomic ./v2/...

coverage.v3:
	$(GO) test -coverpkg=./v3 -coverprofile="$(COVERFILE)" -covermode=atomic ./v3/...

coverage.v3_1:
	$(GO) test -coverpkg=./v3_1 -coverprofile="$(COVERFILE)" -covermode=atomic ./v3_1/...

## -------------------------

## Format & lint
fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

lint: fmt vet

## Clean
clean:
	rm -f "$(COVERFILE)" "$(COVERHTML)"

## Help
help:
	@echo "Targets:"
	@echo "  make deps              - Download dependencies"
	@echo "  make build             - Build all packages"
	@echo "  make test              - Run all tests"
	@echo "  make test.ci           - Run tests with coverage (profile)"
	@echo "  make test.html         - Open coverage report in browser"
	@echo "  make test.func         - Show coverage by function"
	@echo "  make coverage.save     - Save HTML coverage report"
	@echo "  make test.v2           - Run tests only for v2"
	@echo "  make test.v3           - Run tests only for v3"
	@echo "  make test.v3_1         - Run tests only for v3_1"
	@echo "  make coverage.v2       - Coverage report only for v2"
	@echo "  make coverage.v3       - Coverage report only for v3"
	@echo "  make coverage.v3_1     - Coverage report only for v3_1"
	@echo "  make fmt vet lint      - Formatters/Linters"
	@echo "  make clean             - Remove coverage files"
