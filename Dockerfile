FROM golang:1.18.4 as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o bot .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/bot .
COPY --from=builder /app/bot.yml .
COPY --from=builder /app/locales locales

CMD ["./bot"]