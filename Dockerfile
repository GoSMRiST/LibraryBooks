FROM golang:1.25.5-bookworm AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o /cmd/app/exe ./cmd/app/main.go

EXPOSE 8080 50051

CMD ["/cmd/app/exe"]