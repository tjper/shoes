FROM golang:latest
RUN export GOBIN=$GOPATH/bin
RUN echo $GOPATH
RUN echo $GOBIN
RUN echo $PATH

# Create log folder
RUN mkdir /go/log 

RUN go get github.com/lib/pq 

WORKDIR /go/src/github.com/tjper/shoes
COPY . . 
RUN go install 
RUN ls -la

ENTRYPOINT $GOPATH/bin/shoes
EXPOSE 8080
