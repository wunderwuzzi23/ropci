#!/bin/sh

echo "Building"


ver=$(git describe --tags)
echo "Latest tag/version: $ver"


echo "Preparing bin folder"
rm -r ./bin

echo "Creating bin folders for various platforms"

mkdir -p ./bin/windows
mkdir -p ./bin/linux
mkdir -p ./bin/darwin


echo "Building things"

GOOS=linux   GOARCH=amd64 go build -ldflags "-s -w -X 'ropci/cmd.VersionInfo=$(git rev-parse HEAD)'"  -o bin/linux/ropci main.go
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X 'ropci/cmd.VersionInfo=$(git rev-parse HEAD)'" -o bin/windows/ropci.exe main.go
GOOS=darwin  GOARCH=amd64 go build -ldflags "-s -w -X 'ropci/cmd.VersionInfo=$(git rev-parse HEAD)'" -o bin/darwin/ropci main.go

echo "Copying templates"

cp -r ./templates ./bin/windows
cp -r ./templates ./bin/linux
cp -r ./templates ./bin/darwin

echo "Zipping packages"

cd bin
cd windows && zip ../ropci-windows-$ver.zip -r * && cd ..
cd darwin  && zip ../ropci-darwin-$ver.zip -r * && cd ..
cd linux   && zip ../ropci-linux-$ver.zip -r  * && cd ..
cd ..
pwd
echo "Done"

