# Include variables from the .envrc file
include .envrc

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
	@go run ./cmd/api -port=${PORT} -dsn=${DB_DSN} -smtp-sender=${SMTP_SENDER} -smtp-username=${SMTP_USERNAME} -smtp-password=${SMTP_PASS} -smtp-host=${SMTP_HOST} -smtp-port=${SMTP_PORT}

## db/mysql: connect to the database using mysql
.PHONY: db/mysql
db/mysql:
	@mysql -u ${DB_USER} -p${DB_PASSWORD}

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

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api
	GOOS=windows GOARCH=amd64 go build -ldflags='-s' -o=./bin/windows_amd64/api ./cmd/api

# =================================================================================== #
# PRODUCTION
# =================================================================================== #

production_host_ip = "192.168.122.144"

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	@sshpass -p ${SSH_PASS} ssh manager@${production_host_ip}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api threadive@${production_host_ip}:~
	rsync -rP --delete ./migrations threadive@${production_host_ip}:~
	rsync -P ./remote/production/api.service threadive@${production_host_ip}:~
	rsync -P ./remote/production/Caddyfile threadive@${production_host_ip}:~
	ssh -t threadive@${production_host_ip} '\
		migrate -path ~/migrations -database $$FORUM_DB_DSN up \
		&& sudo mv ~/api.service /etc/systemd/system/ \
		&& sudo systemctl enable api \
		&& sudo systemctl restart api \
		&& sudo mv ~Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy \
		'

## bin/api: execute the bin/api application in ./bin/linux_amd64/api
.PHONY: bin/api
bin/api:
	@echo 'Executing binary...'
	@./bin/linux_amd64/api -port=${PORT} -dsn=${DB_DSN} -smtp-username=${SMTP_USERNAME} -smtp-password=${SMTP_PASS} -smtp-host=${SMTP_HOST} -smtp-port=${SMTP_PORT}