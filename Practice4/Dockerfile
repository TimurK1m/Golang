FROM golang:1.25.6-alpine

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект (включая папку database/migrations)
COPY . .

RUN go build -o myapp ./cmd/api/main.go

# Устанавливаем переменную по умолчанию, чтобы работало на Windows/Mac
ENV DB_HOST=host.docker.internal
ENV APP_PORT=8000

CMD [ "./myapp" ]