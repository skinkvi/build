# Используем официальный образ Nginx
FROM nginx:alpine

# Копируем ваши статические файлы в директорию Nginx
COPY . /usr/share/nginx/html

# Копируем конфигурацию Nginx
COPY default.conf /etc/nginx/conf.d/default.conf

# Открываем порт 80
EXPOSE 80

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]