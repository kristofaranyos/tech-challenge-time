#!/bin/sh

if [ $# -ne 1 ]
then
	echo "Missing argument: image name"
	exit 1
fi

docker build . -t $1
docker push $1