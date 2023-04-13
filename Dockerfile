FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /usr/bin/go-server-sample

CMD ["/usr/bin/go-server-sample"]