FROM golang:1.19-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o ./capturer

# Build runtime
FROM alpine:3.8 as runtime

COPY --from=builder /app/capturer /capturer
CMD ["/capturer"]