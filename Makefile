build:
	go build -o telegram-bot

build-rpi:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o telegram-bot-rpi