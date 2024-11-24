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