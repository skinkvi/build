# Используем официальный образ Golang в качестве базового
FROM golang:1.23.2-alpine

# Устанавливаем необходимые зависимости и инструменты, включая OpenVPN
RUN apk update && apk add --no-cache git bash curl build-base python3 py3-pip postgresql-client

# Устанавливаем migrate через go install
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
    mv /go/bin/migrate /usr/local/bin/

RUN pip3 install --no-cache-dir --break-system-packages openai requests

# Копируем файлы приложения
COPY . /app

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы приложения
COPY . .

# Компилируем Go-приложение
RUN go build -o podbor ./cmd/podbor

# Копируем скрипт для миграций
COPY migrate.sh /migrate.sh
RUN chmod +x /migrate.sh

# Копируем конфигурационный файл Tern
COPY tern.conf /app/tern.conf

COPY .env /app/.env
ENV ENV_FILE=/app/.env

# Указываем порт, на котором будет работать приложение
EXPOSE 8080

# Определяем команду для запуска миграций и бэкенда
CMD ["/migrate.sh"]
