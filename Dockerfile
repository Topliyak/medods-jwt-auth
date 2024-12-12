FROM golang:1.23.4-alpine3.21

WORKDIR /medods-app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./bin/app

CMD ["./bin/app"]
