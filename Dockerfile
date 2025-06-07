FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates make bash git
COPY . .
RUN go mod tidy
RUN go build -o bin/cotify-linux-amd64 ./cmd/server
RUN chmod +x /app/bin/cotify-linux-amd64

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache \
    ca-certificates \
    tini \
    tzdata

ENV DB_HOST=localhost
ENV DB_PORT=3306
ENV DB_USER=cotify
ENV DB_PASSWORD=cotify
ENV DB_NAME=cotify
ENV SERVER_PORT=8080
ENV TZ=Asia/Shanghai

COPY --from=builder /app/bin/cotify-linux-amd64 /usr/local/bin/cotify

# install
ENTRYPOINT ["tini", "--"]

CMD ["cotify"]
