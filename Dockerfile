FROM golang:1.16-alpine

WORKDIR /go/src/github.com/dompedro/rpsls
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["rpsls"]
