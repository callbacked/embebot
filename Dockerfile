FROM golang:1.22

WORKDIR /app

COPY main.go .
COPY go.mod .
COPY go.sum .
COPY suppression.go .
COPY match.json .
COPY config.ini .

RUN go mod download
RUN go build -o main .

CMD ["./main"]