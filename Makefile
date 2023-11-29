include .envrc

.PHONY: run run/live test audit tidy db/migrate_up db/migrate_down db/migrate_force db/migrate_create docker/build docker/run

MAIN_PACKAGE_PATH := ./cmd/api
BINARY_NAME := lighter-pics

build:
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

buildWebhook:
	@go build -o ./cmd/setTelegramWebhookUrl ./cmd/setTelegramWebhookUrl.go

run:
	@go run ./cmd/api/main.go -db-dsn=${DB_DSN} -pusher-app-id=${PUSHER_APP_ID} -pusher-key=${PUSHER_KEY} -pusher-secret=${PUSHER_SECRET} -pusher-cluster=${PUSHER_CLUSTER}

test:
	@go test -v ./...

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

db/migrate_up:
	migrate -path=./migrations -database=${DB_DSN} up

db/migrate_down:
	migrate -path=./migrations -database=${DB_DSN} down $(n)

db/migrate_force:
	migrate -path=./migrations -database=${DB_DSN} force $(version)

db/migrate_create:
	migrate create -seq -ext=.sql -dir=./migrations $(name)

docker/build:
	docker build -t imagewizard .

docker/run:
	docker run -e CONFIG_PATH="/app/config.yaml" -e DB_DSN=${DB_DSN} -e ENV="production" -e PUSHER_APP_ID=${PUSHER_APP_ID} -e PUSHER_KEY=${PUSHER_KEY} -e PUSHER_SECRET=${PUSHER_SECRET} -e PUSHER_CLUSTER=${PUSHER_CLUSTER} -e PORT=3001 -e UPLOAD_DIR="/app/uploads" -p 3001:3001 imagewizard


