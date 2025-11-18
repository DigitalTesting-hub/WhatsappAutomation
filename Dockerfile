FROM golang:1.21-alpine

# Install required system dependencies
RUN apk add --no-cache \
    gcc \
    musl-dev \
    sqlite-dev

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o main .

EXPOSE 3000

CMD ["./main"]
