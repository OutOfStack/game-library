# build
FROM golang:1.22-alpine3.19 as builder

WORKDIR /tmp/game-library-api

# copy and download dependencies
COPY go.mod go.sum  ./
RUN go mod download

# copy config and code into container
COPY ./app.env ./out/
COPY . .

# build app
RUN go build -o ./out/game-library-api cmd/game-library-api/main.go

# run
FROM alpine:3.19

WORKDIR /app

# copy built app into runnable container
COPY --from=builder /tmp/game-library-api/out ./

EXPOSE 8000

ENTRYPOINT ["./game-library-api"]
