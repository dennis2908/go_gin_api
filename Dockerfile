FROM golang:1.20 as build

WORKDIR /usr/local/go/src/go-restapi-gin

COPY . . 

RUN go mod download

RUN go mod tidy

RUN go build -o main .


EXPOSE 8080

CMD ["go", "run", "main.go"]