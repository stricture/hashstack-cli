#!/bin/bash

build_dir="hashstack-cli_$1"
mkdir -p $build_dir/{linux,osx,windows}

echo "Updating dependencies"
go get -u
go get github.com/inconshreveable/mousetrap

echo "Building Linux"
go build -o hashstack
mv hashstack "$build_dir/linux/"

echo "Building Windows"
GOARCH=amd64 GOOS=windows go build -o hashstack.exe
mv hashstack.exe "$build_dir/windows/"

echo "Building OS X"
GOARCH=amd64 GOOS=darwin go build -o hashstack
mv hashstack "$build_dir/osx/"

echo "Generating man pages"
cd doc_gen
go get -u
mkdir man
go build
./doc_gen
rm doc_gen
cd ../
mv doc_gen/man "$build_dir/"

7z a "$build_dir.7z" $build_dir
rm -rf $build_dir
