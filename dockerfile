# Используем официальный образ Golang в качестве базового образа
FROM golang:1.20-alpine

# Устанавливаем зависимости для компиляции проекта
RUN apk add --no-cache git

# Копируем исходный код проекта
COPY . /go_sitestat
WORKDIR /go_sitestat

# Загружаем необходимые модули
RUN go mod download

# Компилируем проект
RUN go build -o ./bin/go_sitestat /go_sitestat/main.go

# Настраиваем переменные окружения для подключения к базам данных
ENV REDIS_HOST redis
ENV REDIS_PORT 6379
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_HOST postgres
ENV POSTGRES_PORT 5432
ENV POSTGRES_DATABASE testdb

# Устанавливаем зависимости для работы с Redis и PostgreSQL
RUN apk add --no-cache redis postgresql-client

# Команда запуска приложения
CMD ["./bin/go_sitestat"]