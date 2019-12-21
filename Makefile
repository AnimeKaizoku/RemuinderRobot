format:
	goimports -w .

test:
	rm -f integration_test.db
	TEST_DB_FILE=integration_test.db \
	go test ./...

build:
	go build -o telegram-bot

generate:
	go generate ./...
	goimports -w .

build-rpi:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o telegram-bot-rpi