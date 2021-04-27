FROM golang:1.16.3-buster

WORKDIR /src

COPY . .

RUN ["go", "mod", "download"]

RUN ["go", "mod", "verify"]

RUN ["go", "build", "-o", "build/apiserver", "cmd/apiserver/main.go"]

ENTRYPOINT [ "build/apiserver",  "-debug=true" ]