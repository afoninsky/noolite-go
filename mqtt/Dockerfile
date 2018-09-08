FROM golang:1.11 AS builder
WORKDIR /
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /noolite .

FROM scratch
COPY --from=builder /noolite ./
ENTRYPOINT ["./noolite"]