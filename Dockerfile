FROM golang:1.12.6-alpine3.9

ENV PORT 8080
EXPOSE 8080

RUN mkdir /app
ADD . /app
WORKDIR  /app

ENV GOPROXY https://goproxy.io
ENV GIN_MODE release
RUN go build -tags=jsoniter -o gin-test .

CMD ["/app/gin-test"]
