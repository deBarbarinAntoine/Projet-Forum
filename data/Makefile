# =================================================================================== #
# HELPERS
# =================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# =================================================================================== #
# DEVELOPMENT
# =================================================================================== #

## run/data: run the cmd application
.PHONY: run/data
run/data:
	@go run ./cmd

# =================================================================================== #
# QUALITY CONTROL
# =================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies'
	go mod vendor

# =================================================================================== #
# BUILD
# =================================================================================== #

## build/data: build the cmd application
.PHONY: build/data
build/data:
	@echo 'Building cmd...'
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/data ./cmd
	GOOS=windows GOARCH=amd64 go build -ldflags='-s' -o=./bin/windows_amd64/data ./cmd
