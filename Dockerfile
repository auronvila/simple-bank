#Build  stage
FROM golang:1.23-alpine3.21 AS builder

WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

#Run stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY start.sh .
COPY wait-for.sh .
#COPY app.env .
COPY db/migration ./migration

EXPOSE 3002
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
