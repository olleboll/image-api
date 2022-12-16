FROM golang:1.19

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build cmd/main.go

EXPOSE 8000