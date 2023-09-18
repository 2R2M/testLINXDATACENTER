
FROM golang:latest


WORKDIR /app

COPY . .
# Собираем Go приложение
RUN go build -o app

CMD ["./app"]
