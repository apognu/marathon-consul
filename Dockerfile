FROM golang:alpine
MAINTAINER Antoine POPINEAU <antoine.popineau@appscho.com>

RUN apk update && apk add git

WORKDIR /go/src/github.com/apognu/marathon-consul
COPY . /go/src/github.com/apognu/marathon-consul
RUN go get ./... && go install .

ENTRYPOINT [ "/go/bin/marathon-consul" ]
