FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rest-app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/rest-app .
RUN chmod +x /app/rest-app
RUN ls -la /app


ENV PORT=8888

CMD ["./rest-app"]