FROM golang:1.24

WORKDIR /usr/app

COPY main.go ./
COPY go.mod ./

EXPOSE 8080

RUN go build -v -o /usr/local/bin/

CMD ["/usr/local/bin/hellogo"]