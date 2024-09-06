# Базовый образ Ubuntu
FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    golang \
    git \
    build-essential \
    sqlite3

# Создаем директорию для приложения
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/main cmd/main.go

# Открываем порт
EXPOSE ${TODO_PORT}

ENTRYPOINT ["/app/main"]