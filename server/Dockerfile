FROM golang:1.15.2

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

ADD ./insta-back/server /go/src

WORKDIR /go/src/instagram

RUN export GO111MODULE=on &&\
    go get github.com/go-sql-driver/mysql &&\
    go get github.com/beego/bee &&\
    go get github.com/astaxie/beego

CMD bee run