FROM golang:alpine AS builder

RUN apk add --update gcc musl-dev

ENV GO111MODULE=on
WORKDIR /app

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s"

FROM alpine

RUN apk add --no-cache sqlite

COPY --from=builder /app/turnik-bot /app/turnik-bot

VOLUME [ "/app/data" ]

ENV DB_TYPE sqlite3
ENV DB_URI /app/data/turnik.db
ENV TELEGRAM_URL https://api.telegram.org
ENV TELEGRAM_TOKEN YOUR_TOKEN_FROM_BOTFATHER

ENTRYPOINT ["/app/turnik-bot"]
