FROM golang:alpine

WORKDIR /build

COPY go.mod go.sum ./
COPY *.go ./
COPY entrypoint.sh /

RUN go mod download

RUN go build -o /app

ENTRYPOINT ["/entrypoint.sh"]
