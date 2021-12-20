FROM golang:latest

WORKDIR $GOPATH/src/github.com/drmendoz/iglesias-backend

COPY . . 

RUN go build -o main .

EXPOSE 8000

CMD ["./main"]

