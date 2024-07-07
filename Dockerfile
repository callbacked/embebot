FROM golang:1.22

WORKDIR /app

COPY main.go .
COPY suppression.go .
COPY commands.go .
COPY message.go .
COPY go.mod .
COPY go.sum .
COPY match.json .
COPY config.ini .

RUN go mod download
RUN go build -o main .

CMD ["./main"]