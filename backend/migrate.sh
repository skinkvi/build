#!/bin/bash
set -e

# Загрузка переменных окружения из .env, если они доступны
if [ -f "$ENV_FILE" ]; then
  export $(grep -v '^#' "$ENV_FILE" | xargs)
fi

# Экспортируем пароль для psql
export PGPASSWORD="$DB_PASSWORD"

# Функция проверки готовности PostgreSQL
check_postgres() {
  psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c '\q' >/dev/null 2>&1
}

# Ждем, пока PostgreSQL будет готов
echo "Ожидание доступности PostgreSQL..."
until check_postgres; do
  echo "PostgreSQL недоступен - ждем..."
  sleep 2
done
echo "PostgreSQL готов - применение миграций"

# Применяем миграции
for file in ./migrations/*.sql; do
    echo "Применение миграции $file"
    psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "$file"
done

echo "Миграции успешно применены. Запуск приложения..."

# Запускаем Go-приложение
./podbor