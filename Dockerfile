FROM golang:1.21.6 AS builder
WORKDIR /go/src/github.com/rinchsan/dakoku
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o /go/bin/dakoku ./main.go

FROM chromedp/headless-shell:stable
WORKDIR /go/src/github.com/rinchsan/dakoku
COPY --from=builder /go/bin/dakoku /go/bin/dakoku
ENTRYPOINT ["/go/bin/dakoku"]
