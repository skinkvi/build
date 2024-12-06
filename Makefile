# запуск программы
run:
	docker-compose up --build

# удалить все контейры 
down:
	docker-compose down --rmi all --volumes --remove-orphans
