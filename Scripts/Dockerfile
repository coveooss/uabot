FROM golang

#get source code
RUN git clone https://github.com/coveo/uabot.git
WORKDIR /go/uabot
RUN go get -d

EXPOSE 8080:8080

#run server
CMD [ "go", "run", "server.go", "-queue-length=20", "-port=8080" ]
