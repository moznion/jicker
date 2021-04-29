PKGS := $(shell go list ./...)

check: fmt-check test lint vet
check-ci: fmt-check test-ci vet

test:
	go test -v -cover -race $(PKGS)

test-ci:
	go test -v -cover -race -coverprofile=coverage.txt -covermode=atomic $(PKGS)

lint:
	golangci-lint run -v

vet:
	go vet $(PKGS)

fmt-check:
	goimports -l *.go **/*.go | grep [^*][.]go$$; \
	EXIT_CODE=$$?; \
	if [ $$EXIT_CODE -eq 0 ]; then exit 1; fi \

fmt:
	gofmt -w -s *.go **/*.go
	goimports -w *.go **/*.go

lint-fix:
	golangci-lint run -v --fix

fix:
	$(MAKE) fmt
	$(MAKE) lint-fix
