FROM golang:1.19-alpine3.18 as build

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY . .

RUN go mod tidy
RUN mkdir build
RUN go build -o ./build/server ./cmd/main.go

# ---- backend server ----

FROM alpine:3.18

WORKDIR /usr/src/app

COPY --from=build /usr/src/app/build/server /usr/src/app/server

RUN chmod +x /usr/src/app/server

EXPOSE 8080

CMD ["/usr/src/app/server"]