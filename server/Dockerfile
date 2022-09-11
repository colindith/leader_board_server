FROM golang:1.19.1-buster

WORKDIR /usr/src/app
#RUN apt update && \
#    apt install vim gettext -y \
#    unzip
#COPY ./go.* ./

COPY go.mod .
COPY go.sum .

RUN go env -w GO111MODULE=off && go mod download

COPY . /usr/src/app/.

# Update the repository
# "/usr/src/app/" is added in the directory to make sure that will always be in that location
#COPY . ./

ENV GO111MODULE on

#RUN ["chmod", "+x", "/usr/src/app/dev-entrypoint.sh"]
#RUN ls /usr/src/app/
#RUN ls /usr/src/app/dev-entrypoint.sh
#ENTRYPOINT ["./dev-entrypoint.sh"]
#CMD ["go", "run", "/usr/src/app/main/main.go"]
ENTRYPOINT ["/usr/src/app/dev-entrypoint.sh"]

EXPOSE 8080