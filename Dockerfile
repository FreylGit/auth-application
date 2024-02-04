# Используем базовый образ Go для сборки приложения
FROM golang:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы необходимые для сборки и запуска сервера
COPY ./ /app

# Задаем переменные окружения, в данном случае путь к конфигурационному файлу
ENV CONFIG_PATH=/app/config/dev.env

# Устанавливаем зависимости
RUN go mod download

# Команда сборки приложения
RUN go build -o auth-application ./cmd/auth-application

# Определяем команду для запуска сервера при запуске контейнера
CMD ["./auth-application", "-config=./config/dev.env"] && ["make", "migrateupdev"]
