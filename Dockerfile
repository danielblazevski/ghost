From golang:1.8

RUN mkdir -p $GOPATH/src/go-storage
RUN mkdir -p /base/

WORKDIR $GOPATH/src/go-storage
ADD . $GOPATH/src/go-storage
RUN  go install ./cmd/...
#CMD ["bash"]
