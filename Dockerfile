FROM golang:1.23-alpine3.19 AS builder

COPY . /musaitHrMgBotGo/
WORKDIR /musaitHrMgBotGo/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /musaitHrMgBotGo/bin/bot/ .
COPY --from=0 /musaitHrMgBotGo/configs/ configs/

EXPOSE 80

CMD ["./bot"]