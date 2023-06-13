FROM golang:1.19.0

RUN apt-get update
RUN apt-get install vim nano -y

WORKDIR /usr/src/app
RUN chmod +x /usr/src/app

RUN go install github.com/cosmtrek/air@latest

COPY . .

RUN go mod tidy