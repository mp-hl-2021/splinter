FROM golang:1.16.2-alpine3.13 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o splinter main.go

FROM python:3.8-slim
COPY --from=builder /build/splinter /splinter
RUN pip install pygments

ENTRYPOINT  [ "/splinter" ]
