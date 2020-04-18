format:
	goimports -w .

test: test-unit test-integration test-e2e

test-integration:
	rm -f integration_test.db
	TEST_DB_FILE=integration_test.db \
	go test ./...

test-e2e:
	rm -f e2e/e2e_test.db
	TEST_E2E_DB_FILE=e2e_test.db \
	go test -v e2e/e2e_test.go

test-unit:
	go test ./...

build:
	go build -o telegram-bot

mocks:
	go generate ./...
	goimports -w .

build-rpi:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o telegram-bot-rpi