# ALPINE: Build executable binary
FROM golang:alpine AS builder

# Install Git
RUN apk update && apk add --no-cache git

WORKDIR /go/src/unicorn
COPY . .

RUN go get github.com/gin-gonic/gin
#RUN go get github.com/go-playground/validator

RUN go get -d -v

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/unicorn

# SCRATCH: build small image
FROM alpine:latest
COPY --from=builder /go/bin/unicorn /go/bin/unicorn

WORKDIR go/src/unicorn

COPY templates ./templates

EXPOSE 8080

ENTRYPOINT ["/go/bin/unicorn"]
