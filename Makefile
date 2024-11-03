.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go
run: build
	./.bin/bot

build-image:
	docker build -t musait_hr_bot:v0.1 .

start-container:
	docker run --name musait_hr_bot -p 80:80 --env-file .env musait_hr_bot:v0.1