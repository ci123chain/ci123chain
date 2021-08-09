#! /bin/bash
mkdir -p /go/pkg/mod/github.com/ci123chian
cp -r ./wasmer-go@v1.0.3-rc2 /go/pkg/mod/github.com/ci123chain
echo 'export GOPATH=/go' >> /root/.bashrc
source /root/.bashrc

