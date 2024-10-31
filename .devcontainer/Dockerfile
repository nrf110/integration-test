FROM mcr.microsoft.com/devcontainers/go:1.23-bookworm

ENV GO111MODULE=on
ENV GOBIN=/usr/local/bin

ARG PREFIX="/usr/local"
ARG VERSION="1.39.0"

USER root

RUN chown -R vscode:vscode /usr/local

USER vscode

RUN curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m).tar.gz" | \
    tar -xvzf - -C "${PREFIX}" --strip-components 1

RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
RUN go install github.com/onsi/ginkgo/v2/ginkgo@latest
