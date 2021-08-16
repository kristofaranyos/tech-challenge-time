#!/bin/sh

if [ $# -ne 1 ]
then
	echo "Missing argument: image name"
	exit 1
fi

export IMAGE_NAME=$1
cat k3sconfig.yaml | envsubst | k3s kubectl apply -f -