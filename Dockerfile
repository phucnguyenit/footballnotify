FROM golang:1.12-stretch
RUN mkdir -p "$GOPATH/src/github.com/thang14/footballnotify"
WORKDIR /go/src/github.com/thang14/footballnotify
RUN apt-get update
ADD . .
#RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
#RUN dep ensure
WORKDIR /go/src/github.com/thang14/footballnotify/cmd/footballnotify
RUN go install
WORKDIR /go/bin
ENTRYPOINT ["./footballnotify"]
