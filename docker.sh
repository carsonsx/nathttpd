#!/usr/bin/env bash

# mkdir ~/nathttpd in docker machine
# copy Dockerfile, nathttpd, docker.sh to ~/nathttpd
# run ./docker.sh

docker stop nathttpd
docker rm nathttpd
docker build -t carsonsx/nathttpd .
docker run -itd --name nathttpd --restart=always carsonsx/nathttpd -u amqp://0.0.0.0:5672