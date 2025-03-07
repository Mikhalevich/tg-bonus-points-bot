FROM golang:1.23-alpine3.20 AS builder

WORKDIR /app

RUN GOBIN=/app go install github.com/rubenv/sql-migrate/...@v1.6.1

FROM alpine:3.20

WORKDIR /app/

COPY --from=builder /app/sql-migrate /app/sql-migrate
COPY config/dbconfig-example.yml /app/dbconfig-example.yml
COPY script/db/migrations /app/script/db/migrations

ENTRYPOINT ["./sql-migrate", "up", "-config", "dbconfig-example.yml"]
