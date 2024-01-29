FROM golang:1.21.6 as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o bot .

FROM alpine:3.19
RUN apk update --no-cache && apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bot .
COPY --from=builder /app/bot.yml .
COPY --from=builder /app/locales locales

CMD ["./bot"]