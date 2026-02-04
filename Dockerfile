# Dockerfile
FROM golang:1.25-alpine AS builder
RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go mod tidy
RUN go build -o demo .


FROM busybox:latest
COPY --from=builder /src/demo /demo
EXPOSE 3000
ENTRYPOINT [ "/demo" ]