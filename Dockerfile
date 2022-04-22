FROM golang:1.18 as builder
ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix nocgo -o main .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN apk add --no-cache tzdata

CMD ["./main"]