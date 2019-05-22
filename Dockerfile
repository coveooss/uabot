FROM golang:1.11

WORKDIR $GOPATH/src/github.com/coveo/uabot

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

CMD [ "uabot" ]