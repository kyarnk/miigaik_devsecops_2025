# Этап сборки Go-приложения
FROM docker.io/golang:1.21-alpine AS builder

# Установка рабочей директории
WORKDIR /app

# Копируем файлы go.mod и go.sum (если есть)
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Этап запуска (остается без изменений)
FROM docker.io/alpine:latest

RUN apk --no-cache add ca-certificates
RUN mkdir -p /app/static
COPY --from=builder /app/main /app/main
COPY static /app/static

EXPOSE 8080
WORKDIR /app
CMD ["./main"]