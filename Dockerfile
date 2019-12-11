FROM golang:1.13 as server
WORKDIR /code
ENV GOOS=linux
ENV GARCH=amd64
ENV CGO_ENABLED=0
ENV GO111MODULE=on
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./main.go ./main.go
COPY ./pkg ./pkg
RUN go build -o server

FROM alpine:3.10
RUN apk add iptables
# RUN apk add wireguard-tools
COPY --from=server /code/server /server
CMD /server
