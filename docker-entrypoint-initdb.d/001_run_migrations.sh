#!/bin/bash
set -e

echo "Запуск миграций tern..."

# Ожидание запуска PostgreSQL
until pg_isready -h localhost -p 5432 -U "$POSTGRES_USER"; do
  echo "Ожидание запуска PostgreSQL..."
  sleep 1
done

# Запуск миграций tern
tern migrate -c /migrations/tern.conf

echo "Миграции завершены." 