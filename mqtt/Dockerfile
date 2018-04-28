FROM golang:1.9 as builder

WORKDIR /go/src/github.com/afoninsky/noolite-go/mqtt
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo -o server

FROM alpine:latest
COPY --from=builder /go/src/github.com/afoninsky/noolite-go/mqtt/server /server
CMD ["/server"]