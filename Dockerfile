#Build  stage
FROM golang:1.23-alpine3.21 AS builder

WORKDIR /app
COPY . .
RUN go build -o main main.go

#Run stage
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY start.sh .
COPY wait-for.sh .
#COPY app.env .
COPY db/migration ./db/migration

EXPOSE 3002
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
