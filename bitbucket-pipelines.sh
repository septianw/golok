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
go test -cover -v -coverprofile=coverage.txt -covermode=atomic
bash <(curl -s https://codecov.io/bash) -t 74f70a61-9d13-4257-84b6-6aa0237b2b6b
if [[ $? -eq 0 ]]
then
  GOPATH=$OGOPATH
  exit 0
fi
