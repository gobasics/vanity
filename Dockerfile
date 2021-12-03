FROM golang:1.17 as builder
WORKDIR /
ENV GO111MODULE on
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vanity .

FROM scratch
WORKDIR /
COPY --from=builder /vanity /vanity
ENTRYPOINT ["/vanity"]
