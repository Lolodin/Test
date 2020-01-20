FROM golang:1.13.6 as build
WORKDIR /go/src
COPY ./ ./
RUN go mod download
RUN go build
EXPOSE 8080
CMD [ "./Test" ]
