FROM golang:1.11 AS builder
# https://golang.org/doc/install/source#environment
# [raspberry] uname -a => armv7l => arm
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH arm
WORKDIR /
COPY . ./
RUN go build -a -installsuffix nocgo -o /noolite .

FROM scratch
COPY --from=builder /noolite ./
ENTRYPOINT ["./noolite"]