FROM golang:1.10 AS builder

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR $GOPATH/src/github.com/afoninsky/noolite-go
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN cd mqtt && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM scratch
COPY --from=builder /app ./
ENTRYPOINT ["./app"]