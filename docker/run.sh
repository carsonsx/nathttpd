#!/usr/bin/env bash
set -x
#set -e
set -o pipefail
#
# This script is meant for quick & easy run latest image via:
#   'curl -sSL https://raw.githubusercontent.com/carsonsx/nathttpd/master/docker/run.sh | sh'
# or:
#   'wget -qO- https://raw.githubusercontent.com/carsonsx/nathttpd/master/docker/run.sh | sh'
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
result=`docker images | grep carsonsx/nathttpd`
if [ -n "$result" ]
then
	docker rmi carsonsx/nathttpd
fi
docker build -t carsonsx/nathttpd .

# Run
docker run -itd --name nathttpd --restart=always carsonsx/nathttpd -u amqp://0.0.0.0:5672