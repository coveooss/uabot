FROM golang

#get source code
RUN git clone https://github.com/adambbolduc/uabot.git
WORKDIR /go/uabot
RUN go get -d

#little tweak to update uabot functionnality
RUN rm -rf /go/src/github.com/erocheleau/uabot/scenariolib
RUN cp -r /go/uabot/scenariolib/ /go/src/github.com/erocheleau/uabot/

EXPOSE 8080

#run server
CMD [ "go", "run", "server.go", "-queue-length=20", "-port=8080" ]
