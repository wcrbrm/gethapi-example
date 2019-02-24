FROM golang:1.11.1-alpine as build
RUN apk add --no-cache git gcc libc-dev

# deps
RUN go get -u "github.com/satori/go.uuid"
RUN go get -u "github.com/jmoiron/sqlx"
RUN go get -u "github.com/lib/pq"
RUN go get -u "github.com/ethereum/go-ethereum/ethclient"

RUN mkdir -p /go/src/github.com/wcrbrm/gethapi-example/server
WORKDIR /go/src/github.com/wcrbrm/gethapi-example/server

# Uploading source files
COPY main.go ./
COPY blockchain ./blockchain/
COPY database ./database/
COPY tcpserver ./tcpserver/

# Compiling (cannot do CGO_ENABLED=0 because of ethereum RPC deps)
RUN go build -a -installsuffix cgo --ldflags "-s -w" -o /usr/bin/server

FROM alpine:3.9
COPY --from=build /usr/bin/server /root/

# genesis block can be used to initialize accounts
# COPY genesis.json /root/
# ENV GETH_GENESIS_PATH=/root/genesis.json
EXPOSE 9091
WORKDIR /root/
CMD ["./server"]