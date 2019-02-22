FROM golang:1.11.5-alpine as build
RUN mkdir -p /go/src/server/github.com/wcrbrm/gethapi-example
WORKDIR /go/src/github.com/wcrbrm/gethapi-example

RUN go get -u "github.com/ethereum/go-ethereum/ethclient"
RUN go get -u "github.com/satori/go.uuid"
RUN go get -u "github.com/jmoiron/sqlx"
RUN go get -u "github.com/lib/pq"

COPY main.go .
COPY blockchain ./blockchain/
COPY database ./database/
COPY tcpserver ./tcpserver/

RUN CGO_ENABLED=0 go build -a -installsuffix cgo --ldflags "-s -w" -o /usr/bin/server

FROM alpine:3.7
COPY --from=build /usr/bin/server /root/
EXPOSE 3000
WORKDIR /root/
CMD ["./server"]