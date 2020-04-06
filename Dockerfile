FROM golang:1.14

WORKDIR /go/src/github.com/peterzernia/set

COPY go.mod /go/src/github.com/peterzernia/set
COPY go.sum /go/src/github.com/peterzernia/set

RUN go mod download

COPY . /go/src/github.com/peterzernia/set

RUN curl -SL https://github.com/maxcnunes/gaper/releases/download/v1.0.3/gaper_1.0.3_linux_amd64.tar.gz | tar -xvzf - -C "${GOPATH}/bin"

EXPOSE 8002
