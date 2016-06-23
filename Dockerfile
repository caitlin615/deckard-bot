FROM golang:1.5

RUN apt-get update

ENV ROOT_PATH /go/src/github.com/handwritingio/deckard-bot
RUN echo 'export PATH=$PATH:/go/bin' >> ~/.bashrc

WORKDIR $ROOT_PATH

RUN go get github.com/golang/lint/golint

ADD . $ROOT_PATH

RUN go get -v -d ./...
RUN go install -v ./...

CMD ["deckard-bot"]
