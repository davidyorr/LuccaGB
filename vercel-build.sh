#!/bin/bash

# download Go
curl -L -o go.tar.gz https://go.dev/dl/go1.24.3.linux-amd64.tar.gz

# extract it
tar -C . -xzf go.tar.gz

# add local Go binary to PATH
export PATH="$(pwd)/go/bin:$PATH"

# verify installation
go version

# run the package.json build script
cd web
pnpm run build