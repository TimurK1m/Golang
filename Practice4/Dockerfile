FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/api/main.go


FROM alpine

# бинарник
COPY --from=builder /app/myapp /myapp

# миграции
COPY --from=builder /app/database /database

ENV DB_HOST=host.docker.internal
ENV APP_PORT=8000

ENTRYPOINT ["/myapp"]