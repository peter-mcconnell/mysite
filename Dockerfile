FROM golang:1.11 as builder
LABEL maintainer "Peter McConnell<me@petermcconnell.com>"
WORKDIR /go/src/github.com/pemcconnell/mysite
COPY . .
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o mysite

FROM scratch
COPY --from=builder /go/src/github.com/pemcconnell/mysite/mysite /usr/bin/app
COPY --from=builder /go/src/github.com/pemcconnell/mysite/static /var/www/static
COPY --from=builder /go/src/github.com/pemcconnell/mysite/templates /var/www/templates
WORKDIR /var/www/
EXPOSE 8080
ENTRYPOINT ["/usr/bin/app"]
