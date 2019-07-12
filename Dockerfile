FROM golang

RUN go get -u github.com/golang/dep/cmd/dep

#RUN mkdir /go/src/github.com
#RUN mkdir /go/src/github.com/smukhalov
#RUN mkdir /go/src/github.com/smukhalov/test4

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/smukhalov/test4

# Setting up working directory
WORKDIR /go/src/github.com/smukhalov/test4

RUN dep ensure

# Build the taskmanager command inside the container.
RUN cd /go/src/github.com/smukhalov/test4

#RUN go run main.go

# Service listens on port 8080.
EXPOSE 8080
EXPOSE 27017
