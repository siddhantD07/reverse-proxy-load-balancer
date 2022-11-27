FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

RUN mkdir /gserve

COPY /src/gserve /gserve

WORKDIR /gserve

RUN go mod init gserve

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -o gserve ./

FROM scratch

COPY --from=builder /gserve .

EXPOSE 80

ENTRYPOINT ["./gserve"]
