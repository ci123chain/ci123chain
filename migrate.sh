#! /bin/bash
mkdir -p /go/pkg/mod/github.com
cp -r ./wasmerio /go/pkg/mod/github.com
echo 'export GOPATH=/go' >> /root/.bashrc
source /root/.bashrc

