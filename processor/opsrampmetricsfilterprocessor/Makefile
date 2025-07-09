.PHONY: build test lint clean fmt vet

GO := go
GOFMT := gofmt

build:
	$(GO) build ./...

test:
	$(GO) test -race -v ./...

fmt:
	$(GOFMT) -w .

vet:
	$(GO) vet ./...

lint: fmt vet

clean:
	$(GO) clean ./...
	rm -f coverage.out

coverage:
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

help:
	@echo "Available targets:"
	@echo "  build    - Build the processor"
	@echo "  test     - Run tests"
	@echo "  fmt      - Format code"
	@echo "  vet      - Run go vet"
	@echo "  lint     - Run fmt and vet"
	@echo "  clean    - Clean build artifacts"
	@echo "  coverage - Generate test coverage report"
