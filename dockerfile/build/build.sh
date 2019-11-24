#!/bin/sh

echo "Scooping up recent copy of code"

tar -cf plunder.tar ../../cmd ../../pkg ../../main.go ../../Makefile ../../go.mod ../../go.sum

docker build -t plunder:build .

docker run -it --rm plunder:build /go/bin/plunder version

echo ""
echo "Remove the code copy rm plunder.tar"
echo "Remove the local test docker container docker rmi plunder:build"