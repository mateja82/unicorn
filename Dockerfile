FROM golang:latest
#RUN mkdir /unicorn
ADD . /go/src/unicorn
# create a working directory
WORKDIR /go/src/unicorn
COPY . .
#RUN go get -d -v ./...
RUN go get github.com/gin-gonic/gin
RUN go get github.com/go-playground/validator
RUN go get github.com/aws/aws-sdk-go/service/ssm
RUN go get github.com/aws/aws-sdk-go/aws
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o unicorn .
CMD ["./unicorn"]
