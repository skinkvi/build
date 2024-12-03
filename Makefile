.PHONY: build up down restart clean logs

# Собрать образы Docker
build:
	docker-compose build

# Запустить контейнеры в фоновом режиме
up:
	docker-compose up -d

# Остановить и удалить контейнеры
down:
	docker-compose down

# Перезапустить контейнеры с учетом изменений
restart: down build up

# Остановить и удалить контейнеры, удалить тома и сети, очистить неиспользуемые данные Docker
clean:
	docker-compose down -v --remove-orphans
	docker system prune -f
	docker volume prune -f
