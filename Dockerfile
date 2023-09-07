# Choose whatever you want, version >= 1.16
FROM golang:1.21-alpine

WORKDIR /go-gorm

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air"]