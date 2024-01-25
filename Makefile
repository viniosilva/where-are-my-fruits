include .env
export

all:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go get

infra/up:
	docker-compose up --build -d

infra/down:
	docker-compose down --remove-orphans

db/migration:
	migrate -database mysql://user_admin:$(MYSQL_PASSWORD)@tcp\(localhost:3307\)/where-are-my-fruits -path db/migrations up

docker/image:
	docker build -t where-are-my-fruits .

swag:
	swag init

clean-mocks:
	rm ./mocks/*_mocks.go

.PHONY: mocks
mocks: clean-mocks
	go generate ./...

.PHONY: tests
tests:
	go test ./... -cover

tests/cov:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

tests/e2e:
	go test ./tests/... -cover

run:
	go run ./main.go