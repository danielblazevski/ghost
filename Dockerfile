From golang:1.8

RUN mkdir -p $GOPATH/src/ghost

WORKDIR $GOPATH/src/ghost
ADD . $GOPATH/src/ghost
RUN  go install ./cmd/...
