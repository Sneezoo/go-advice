FROM golang:latest

ENV MONGO_HOST="mongo"

RUN go get github.com/derekparker/delve/cmd/dlv
RUN go get github.com/golang/dep/cmd/dep
WORKDIR "/go/src/github.com/Sneezoo/advicery"
ENTRYPOINT "./advicery"

EXPOSE 8080

COPY Gopkg.* ./

COPY *.go ./
COPY advice ./advice
RUN dep ensure
RUN go build -gcflags='-N -l' -tags=jsoniter

COPY www/dist ./www
COPY *.sh ./