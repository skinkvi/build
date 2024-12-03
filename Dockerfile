FROM postgres:16

# Установка wget и tern
RUN apt-get update && \
    apt-get install -y wget && \
    wget -q -O /usr/local/bin/tern https://github.com/jackc/tern/releases/download/v2.2.3/tern_linux_amd64 && \
    chmod +x /usr/local/bin/tern

# Копирование файлов миграций
COPY migrations /migrations

# Копирование скрипта entrypoint
COPY docker-entrypoint-initdb.d /docker-entrypoint-initdb.d

# Установка переменных окружения
ENV TERN_MIGRATIONS=/migrations 

FROM golang:1.20-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum для скачивания зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем приложение
RUN go build -o podbor ./cmd/podbor/main.go

# Создаем финальный образ
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/podbor .

# Копируем конфигурационные файлы
COPY --from=builder /app/configs ./configs

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["./podbor"]