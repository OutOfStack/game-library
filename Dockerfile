# build
FROM golang:1.20-alpine3.17 as builder

WORKDIR /tmp/game-library-api

# copy and download dependencies
COPY go.mod go.sum  ./
RUN go mod download

# copy config and code into container
COPY ./app.env ./out/
COPY . .

RUN go build -o ./out/game-library-api cmd/game-library-api/main.go

# run
FROM ubuntu:23.04

WORKDIR /app

COPY --from=builder /tmp/game-library-api/out ./

EXPOSE 8000

ENTRYPOINT ["./game-library-api"]
