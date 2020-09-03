FROM golang:1.14 as builder
ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o main .

FROM scratch
COPY --from=builder /app/main /app/
ENTRYPOINT ["/app/main"]