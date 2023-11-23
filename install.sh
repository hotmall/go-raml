#!/bin/bash
set -ex
go get -u github.com/tools/godep
go get -u github.com/jteeuwen/go-bindata/...
go get -u github.com/hotmall/go-raml
cd $GOPATH/src/github.com/hotmall/go-raml
sh build.sh

