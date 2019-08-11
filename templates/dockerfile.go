package templates

var Dockerfile = `FROM golang as builder

ENV GO111MODULE=on
WORKDIR /go/src/{{.Config.Package}}

COPY . .
RUN go get ./... 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /tmp/app *.go

FROM jakubknejzlik/wait-for as wait-for

FROM alpine:3.5

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=wait-for /usr/local/bin/wait-for /usr/local/bin/wait-for
COPY --from=builder /tmp/app /usr/local/bin/app

# https://serverfault.com/questions/772227/chmod-not-working-correctly-in-docker
RUN chmod +x /usr/local/bin/app

ENTRYPOINT []
CMD [ "/bin/sh", "-c", "wait-for ${DATABASE_URL} && app"]
`
