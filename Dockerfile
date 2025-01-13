FROM golang:1.23

WORKDIR /go/src/app

COPY . .

RUN go mod download

CMD ["go", "run", "cmd/app/main.go"]
