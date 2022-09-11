FROM golang:1.19.1-buster

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .

RUN go env -w GO111MODULE=off && go mod download

COPY . /usr/src/app/.

ENV GO111MODULE on

ENTRYPOINT ["/usr/src/app/dev-entrypoint.sh"]

EXPOSE 8080