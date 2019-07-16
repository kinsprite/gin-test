# build
FROM golang:1.12.6-alpine3.9 as build

ENV PORT 8080
EXPOSE 8080

RUN mkdir /app
ADD . /app

ENV GOPROXY https://goproxy.io
ENV GIN_MODE release

WORKDIR  /app
RUN go mod vendor
RUN go build -mod=vendor -tags=jsoniter -o gin-test .


# release
FROM alpine:3.9
RUN mkdir /app
COPY --from=build /app/gin-test /app/gin-test

WORKDIR  /app
CMD ["/app/gin-test"]
