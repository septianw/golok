#!/bin/bash

OGOPATH=$GOPATH
PWD=$(pwd)
mkdir {bin,src}
export GOPATH=$PWD
mkdir -p src/bitbucket.org/septianw/golok
mv *.go src/bitbucket.org/septianw/golok
cd src/bitbucket.org/septianw/golok
go get
go build
go test -v -cover
if [[ $? -eq 0 ]]
then
  GOPATH=$OGOPATH
  exit 0
fi
