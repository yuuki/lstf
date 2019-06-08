FROM golang:1.12

RUN mkdir -p /src
WORKDIR /src
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
