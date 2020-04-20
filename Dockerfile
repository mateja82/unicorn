FROM golang:latest
RUN mkdir /unicorn
ADD . /unicorn/
# create a working directory
WORKDIR /unicorn
COPY . .
# run main.go
CMD ["go", "run", "unicorn"]
