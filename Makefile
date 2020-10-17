all: main deploy

main: main.go
	GOARCH=arm64 GOARM=7 GOOS=linux go build -o main

deploy:
	scp main pi@192.168.1.45:/home/pi/main
	scp -r ./public pi@192.168.1.45:/home/pi/public/

.PHONY: all deploy