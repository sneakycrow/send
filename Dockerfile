FROM golang:1.15

ENV GO11MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/send /app/cmd/send
EXPOSE 8080
ENTRYPOINT ["/app/send"]
