FROM golang:1.25-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server cmd/main.go

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /build/server ./server
# COPY --from=builder /build/.env.production .env

RUN chmod +x server

EXPOSE 8080

CMD [ "./server" ]