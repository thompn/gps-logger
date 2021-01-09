FROM golang:1.15-alpine
RUN mkdir /api
ADD . /api
WORKDIR /api
RUN go build -o main .
CMD ["/app/main"]