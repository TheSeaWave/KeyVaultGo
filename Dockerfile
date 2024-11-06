# Stage 1: Сборка приложения
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Копируем go.mod и go.sum, затем устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

RUN ls -al /app

# Копируем остальные файлы и собираем приложение
COPY . . 
RUN go build -o server cmd/main.go

# Stage 2: Подготовка образа для запуска
FROM alpine:latest
WORKDIR /app

# Копируем скомпилированный бинарник и конфиги
COPY --from=builder /app/server .
COPY --from=builder /app/db.json /app/db.json


# Запускаем приложение
CMD ["./server"]
