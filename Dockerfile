FROM golang:1.23-alpine3.19 AS builder

COPY . /github.com/mirshodNasilloyev/musait_hr_bot/
WORKDIR /github.com/mirshodNasilloyev/musait_hr_bot/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/mirshodNasilloyev/musait_hr_bot/bin/bot/ .
COPY --from=0 /github.com/mirshodNasilloyev/musait_hr_bot/configs/ configs/

EXPOSE 80

CMD ["./bot"]