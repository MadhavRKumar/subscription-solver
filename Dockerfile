FROM golang:1.22.5 AS dev

WORKDIR /usr/src/app

RUN go install github.com/air-verse/air@v1.52.3

COPY . .
RUN go mod tidy

CMD ["air"]

FROM golang:1.22.5 AS build

WORKDIR /usr/src/app

COPY . .

RUN go mod download

RUN go build -C ./cmd/api-server -o server

FROM golang:1.22.5 AS prod

WORKDIR /usr/src/prod

COPY --from=build /usr/src/app/cmd/api-server/server .

CMD ["./server"]


