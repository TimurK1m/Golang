FROM golang:1.25.6-alpine

WORKDIR /app

COPY . .

RUN go build -o myapp ./cmd/api/main.go
CMD [ "./myapp" ]