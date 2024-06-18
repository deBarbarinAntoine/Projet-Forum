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

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@cd api/ && make run/api

## db/mysql: connect to the database using mysql
.PHONY: db/mysql
db/mysql:
	@cd api/ && make db/mysql

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	cd api/ && make db/migrations/new

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@cd api/ && make db/migrations/up

## db/migrations/drop: drop database migrations
.PHONY: db/migrations/drop
db/migrations/drop: confirm
	@cd api/ && make db/migrations/drop

# =================================================================================== #
# QUALITY CONTROL
# =================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@cd api/ && make audit
	@cd backend/ && make audit

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@cd api/ && make vendor
	@cd backend/ && make vendor

# =================================================================================== #
# BUILD
# =================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@cd api/ && make build/api
	@cd backend/ && make build/backend

# =================================================================================== #
# PRODUCTION
# =================================================================================== #

production_host_ip = "192.168.122.144"

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	@cd api/ && make prodution/connect

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	cd api/ && make production/deploy/api

## production/deploy/backend: deploy the backend to production
.PHONY: production/deploy/backend
production/deploy/backend:
	cd backend/ && make production/deploy/backend

## bin/api: execute the bin/api application in ./bin/linux_amd64/api
.PHONY: bin/api
bin/api:
	@cd api/ && make bin/api

## bin/backend: execute the bin/backend application in ./bin/linux_amd64/backend
.PHONY: bin/backend
bin/backend:
	@cd backend/ && make bin/backend