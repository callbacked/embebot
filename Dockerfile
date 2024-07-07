FROM golang:1.22

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
COPY suppression.go .
COPY commands.go .
COPY message.go .

COPY match.json .
COPY config.ini .


RUN go build -o main .

CMD ["./main"]