FROM golang:1.20-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /dueit/mail-service

COPY . .

RUN go mod download
RUN go build -o mail-service .

FROM alpine:latest

RUN mkdir app \
    && mkdir app/internal \
    && mkdir app/internal/template \
    && mkdir app/internal/template/html

COPY --from=builder /dueit/mail-service/.env /app
COPY --from=builder /dueit/mail-service/internal/template/html /app/internal/template/html
COPY --from=builder /dueit/mail-service/mail-service /app

WORKDIR /app
EXPOSE 9090
CMD ./mail-service