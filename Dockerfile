FROM golang:1.22.5

WORKDIR /usr/src/app

RUN go install github.com/air-verse/air@v1.52.3

COPY . .
RUN go mod tidy
