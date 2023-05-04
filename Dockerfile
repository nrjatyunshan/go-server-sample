FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

#RUN go build -gcflags="-l" -o /usr/bin/go-server-sample
RUN go build -o /usr/bin/go-server-sample

CMD ["/usr/bin/go-server-sample"]