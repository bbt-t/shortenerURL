FROM golang:1.19.3-bullseye

WORKDIR /app
COPY . .

RUN go build -o main ./cmd/shortener/main.go

EXPOSE 8080

CMD ["./main"]
