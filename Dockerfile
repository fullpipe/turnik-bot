FROM golang AS builder

ENV GO111MODULE=on
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s"

FROM scratch

COPY --from=builder /app/turnik-bot /go/bin/turnik-bot

ENV DB_TYPE mysql
ENV DB_URI mysql://adam:alphabetonly@127.0.0.1/test
ENV TELEGRAM_URL https://api.telegram.org
ENV TELEGRAM_TOKEN YOUR_TOKEN_FROM_BOTFATHER

ENTRYPOINT ["/go/bin/turnik-bot"]
