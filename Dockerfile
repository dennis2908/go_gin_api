FROM golang:1.22 as build

WORKDIR /usr/local/go/src/go-restapi-gin

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

ENTRYPOINT ["go", "run", "main.go"]