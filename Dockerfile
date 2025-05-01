FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY main.go .
RUN go build -v -o main main.go
COPY simulation_data simulation_data

CMD ["/app/main"]