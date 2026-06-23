FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/user-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/user-migrate ./cmd/migrate

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /bin/user-api /app/user-api
COPY --from=builder /bin/user-migrate /app/user-migrate
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["/app/user-api"]
