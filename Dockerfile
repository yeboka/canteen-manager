FROM golang:1.22-alpine AS builder

RUN apk update && apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY config/apiserver.toml /app/config/apiserver.toml

RUN make build

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/apiserver .
COPY --from=builder /app/config/apiserver.toml ./config/apiserver.toml

EXPOSE 8080

CMD ./apiserver
