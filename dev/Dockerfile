FROM golang:1.22.2

RUN go install golang.org/x/tools/cmd/goimports@v0.21.0

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.57.2

WORKDIR /go/src/app

RUN go version && golangci-lint --version
