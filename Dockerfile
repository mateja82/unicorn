FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/unicorn
COPY . .

RUN go get github.com/gin-gonic/gin && \
    go get github.com/go-playground/validator && \
    go get -d -v


RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/unicorn

FROM alpine:latest

COPY --from=builder /go/bin/unicorn /go/bin/unicorn
WORKDIR $GOPATH/src/unicorn

COPY templates ./templates
EXPOSE 8080

ENTRYPOINT ["/go/bin/unicorn"]
