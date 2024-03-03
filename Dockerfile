FROM golang:1.22.0-alpine as builder

WORKDIR /app

COPY ./ ./
RUN mkdir bin
RUN go mod tidy
RUN go build -o ./bin/cobaApp ./main.go

FROM alpine:3

WORKDIR /app

COPY --from=builder /app/config.json ./
COPY --from=builder /app/bin/cobaApp ./

EXPOSE 5005
CMD ./cobaApp