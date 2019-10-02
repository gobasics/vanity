FROM golang:1.13 as builder
WORKDIR /
RUN mkdir /etc/ssl/cache
ENV GO111MODULE on
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vanity .

FROM scratch
WORKDIR /
ENV DIR_CACHE /etc/ssl/cache
COPY --from=builder /etc/ssl /etc/ssl
COPY --from=builder /etc/ssl/cache /etc/ssl/cache
COPY --from=builder /vanity /vanity
ENTRYPOINT ["/vanity"]
