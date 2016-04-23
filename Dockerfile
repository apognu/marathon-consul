FROM golang:alpine
MAINTAINER Antoine POPINEAU <antoine.popineau@appscho.com>

RUN apk update && apk add git

WORKDIR /go/src/app
COPY . /go/src/app
RUN go get -d . && go build .

ENTRYPOINT [ "/go/src/app/app" ]
