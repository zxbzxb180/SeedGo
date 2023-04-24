FROM golang:1.18.3-alpine3.15


WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build -o app

CMD [ "./app" ]