exist := $(wildcard .envrc)
ifneq ($(strip $(exist)),)
  include .envrc
endif

.PHONY: run run/live test audit tidy db/migrate_up db/migrate_down db/migrate_force db/migrate_create docker/build docker/run mocks report install test/coverage

MAIN_PACKAGE_PATH := ./cmd/api
BINARY_NAME := cnvrt

build:
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

run:
	@go run ./cmd/api/main.go -db-dsn=${DB_DSN} -metrics-user=${METRICS_USER} -metrics-password=${METRICS_PASSWORD} -firebase-project-id=${FIREBASE_PROJECT_ID} -allow-origins=${ALLOW_ORIGINS} & caddy run

test/coverage:
	@ENV=test go test -v ./... -coverprofile=coverage.out

test:
	@ENV=test go test -v -count=5 ./...

install:
	@go get -u ./...

report:
	@go tool cover -html=coverage.out -o coverage.html

run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
			--build.cmd "make build" --build.bin "/tmp/bin/${BINARY_NAME}" --build.delay "100" \
			--build.exclude_dir "" \
			--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
			--misc.clean_on_exit "true"

audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

tidy:
	go fmt ./...
	go mod tidy -v

mocks:
	mockgen -source=internal/repositories/repositories.go --destination=internal/mocks/repositories.go --package=mocks
	mockgen -source=internal/services/services.go --destination=internal/mocks/services.go --package=mocks
	mockgen -source=internal/types/types.go --destination=internal/mocks/types.go --package=mocks

db/migrate_up:
	echo ${DB_DSN}
	migrate -path=./migrations -database=${DB_DSN} up

db/migrate_down:
	migrate -path=./migrations -database=${DB_DSN} down $(n)

db/migrate_force:
	migrate -path=./migrations -database=${DB_DSN} force $(version)

db/migrate_create:
	migrate create -seq -ext=.sql -dir=./migrations $(name)

docker/build:
	docker build --build-arg db_dsn=${DB_DSN} -t cnvrt .

docker/run:
	docker run -e CONFIG_PATH="/app/config.yaml" -e DB_DSN=${DB_DSN} -e ENV="production" -e METRICS_USER=${METRICS_USER} -e METRICS_PASSWORD=${METRICS_PASSWORD} -e FIREBASE_PROJECT_ID=${FIREBASE_PROJECT_ID} -e PORT=3001 -e UPLOAD_DIR="/app/uploads" -e ALLOW_ORIGINS="https://cnvrt.pics" -p 3001:3001 cnvrt


