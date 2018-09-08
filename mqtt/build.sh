#!/bin/sh
set -eu

TAG="vkfont/pi3-noolite:$1-arm"
echo ">>> $TAG"
docker build -t $TAG .
docker push $TAG
