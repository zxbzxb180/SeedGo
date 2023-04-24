FROM golang:1.18.3-alpine3.15

ENV GOPROXY=https://goproxy.io

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build -o app

CMD [ "./app" ]