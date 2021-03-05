DOCKER_IMAGE_NAME=muktihari/order-transaction-ddd

test:
	go test -race -v -cover ./...
run:
	go run main.go
run-mongo:
	@echo "mongo should have replica(s) to enable transactional"
	go run main.go -repo mongo
run-mongo-migrate:
	@echo "mongo should have replica(s) to enable transactional"
	go run main.go -repo mongo -migrate
build:
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build  -ldflags="-s -w" -trimpath -o app main.go
docker-build:
	docker build -t ${DOCKER_IMAGE_NAME} .