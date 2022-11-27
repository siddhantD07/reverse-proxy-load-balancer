FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

RUN mkdir /proxy

COPY /src/grproxy /proxy

WORKDIR /proxy

RUN go mod init proxy

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -o grproxy ./

FROM scratch

COPY --from=builder /proxy .

EXPOSE 80

ENTRYPOINT ["./grproxy"]
