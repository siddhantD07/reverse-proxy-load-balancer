FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/app/

COPY ./src/grproxy .

RUN go mod init grproxy

RUN go get -d -v

RUN CGO_ENABLED=0 go build -o /go/bin

FROM scratch

COPY --from=builder /go/bin/ .

EXPOSE 80

ENTRYPOINT ["./grproxy"]