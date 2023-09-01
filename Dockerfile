FROM golang:latest
WORKDIR /app
COPY . .
#Скачивание всех зависимостей в mod.go
RUN go mod download
#Установка ENV токена для тг бота
ENV TOKEN=
#Установка ENV порта gRPC сервера
ENV PORT=8080
#Сборка бинаря
RUN go build -o main ./cmd
#Запуск бинаря 
CMD  [ "./main" ] 
EXPOSE 8080