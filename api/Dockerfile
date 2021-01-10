FROM golang:1.15-alpine
RUN mkdir /api
ADD . /api
WORKDIR /api
RUN go build -o main .
COPY ./data/ /docker-entrypoint-initdb.d/
CMD ["/app/main"]