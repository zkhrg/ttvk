# Используем официальный Go-образ на базе Alpine для сборки
FROM ubuntu:latest AS builder

# Устанавливаем зависимости
RUN apt-get update && apt-get install -y iputils-ping ca-certificates golang

# Указываем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum, устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 go build  -o app ./cmd/main.go

# Финальный контейнер с минимальным размером
FROM ubuntu:latest

# Устанавливаем утилиту ping
# RUN apt-get install --no-cache iputils
RUN apt-get update && apt-get install -y iputils-ping ca-certificates

# Копируем бинарник из сборочного контейнера
COPY --from=builder /app/app /app

# Запускаем приложение
CMD ["/app"]