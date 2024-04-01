FROM golang:1.22.0 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main_service

FROM golang:1.22.0

WORKDIR /app

COPY --from=builder /app/. .

CMD ["./my_go_program"]