FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth-migrate ./cmd/migrate

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /bin/auth-api /app/auth-api
COPY --from=builder /bin/auth-migrate /app/auth-migrate

# สำคัญมาก
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["/app/auth-api"]