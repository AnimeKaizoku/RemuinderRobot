format:
	goimports -w .

test: test-unit test-integration

test-integration:
	rm -f integration_test.db
	TEST_DB_FILE=integration_test.db \
	go test ./...

test-unit:
	go test ./...

build:
	go build -o telegram-bot

mocks:
	go generate ./...
	goimports -w .

build-rpi:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o telegram-bot-rpi