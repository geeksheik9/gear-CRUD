FROM golang:alpine AS builder
ARG VERSION
RUN apk add --no-cache --virtual .build-deps git libc6-compat build-base
WORKDIR /gear-CRUD

COPY . .
RUN go mod download
WORKDIR /gear-CRUD/main
RUN go build -gcflags "all=-N -l" -ldflags "-X main.version=${VERSION}" -o app;

FROM alpine:latest
WORKDIR /root
COPY --from=builder /gear-CRUD/main/app .

CMD ["./app"]
