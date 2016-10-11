#!/usr/bin/env bash
set -x
#set -e
set -o pipefail
#
# This script is meant for quick & easy build and push nathttpd image via:
#   'curl -sSL https://raw.githubusercontent.com/carsonsx/nathttpd/master/docker/build.sh | sh'
# or:
#   'wget -qO- https://raw.githubusercontent.com/carsonsx/nathttpd/master/docker/build.sh | sh'
#

# Download Dockerfile
curl -sSLO "https://raw.githubusercontent.com/carsonsx/nathttpd/master/docker/Dockerfile"

# Stop and remove image if exists
result=`docker ps | grep nathttpd`
if [ -n "$result" ]
then
	docker stop nathttpd
fi
result=`docker ps -a | grep nathttpd`
if [ -n "$result" ]
then
	docker rm nathttpd
fi

# Build image
docker build -t carsonsx/nathttpd .
docker tag carsonsx/nathttpd carsonsx/nathttpd:0.1

# Push image
docker push carsonsx/nathttpd
docker push carsonsx/nathttpd:0.1

# Clean
docker rmi carsonsx/nathttpd:0.1
docker rmi carsonsx/nathttpd
rm -rf Dockerfile