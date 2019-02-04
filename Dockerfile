FROM golang:1.11-alpine3.8

ENV GOPATH $HOME/go:$GOPATH
ENV PATH $HOME/go/bin:$PATH

#Change to app folder for go
WORKDIR /go/src/github.com/ravster/pinterest_clone

RUN apk add curl git build-base

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN go get -v "github.com/rubenv/sql-migrate/..."

COPY . .

# RUN dep ensure -v

# RUN go install

# CMD ["go-wrapper", "run"]
