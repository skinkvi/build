#!/bin/sh

# Ожидаем, пока база данных не будет доступна
until nc -z db 5432; do
  echo "Ожидание подключения к базе данных..."
  sleep 1
done

echo "Подключено к базе данных. Применение миграций..."

# Применяем миграции с помощью migrate
migrate -path=/app/migrations -database "postgres://$DB_USER:$DB_PASSWORD@db:$DB_PORT/$DB_NAME?sslmode=disable" up

echo "Миграции успешно применены. Запуск приложения..."

# Запускаем Go приложение
./podbor