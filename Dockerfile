From golang:1.8

ENV GOPATH=/go/ghost
RUN mkdir -p /base/

WORKDIR $GOPATH
ADD . $GOPATH
RUN  go install ./src/cmd/...
#CMD ["bash"]
