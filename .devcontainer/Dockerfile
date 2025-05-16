# Base image for Go
FROM golang:1.23 AS go-base

# Add a non-root user
ARG USER=vscode

RUN apt update \
    && apt install -y --no-install-recommends sudo curl \
    && apt autoremove -y \
    && rm -rf /var/lib/apt/lists/* \
    && useradd -m -s /usr/bin/bash vscode \
    && echo "vscode ALL=(ALL) NOPASSWD: ALL" >/etc/sudoers.d/vscode \
    && chmod 0440 /etc/sudoers.d/vscode

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/air-verse/air@latest

# Switch to vscode user
USER vscode
WORKDIR /home/vscode

# Set Go environment variables
ENV PATH="/usr/local/go/bin:$PATH"
ENV GOPATH="/home/vscode/go"

