FROM golang:1.15-alpine as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a
RUN ls -a ${pwd}/app

FROM scratch
COPY --from=builder /app/send /app/send
EXPOSE 8080
ENTRYPOINT ["/app/send"]
