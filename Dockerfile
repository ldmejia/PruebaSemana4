FROM golang:1.24.3-alpine

WORKDIR /app 

COPY go.mod ./

RUN go mod download 

COPY . . /app/

RUN go build -o app . 

CMD ["./app"]

