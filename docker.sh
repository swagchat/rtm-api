#!/bin/sh

docker build . -t docker.swagchat.io:30000/rtm-api && docker push docker.swagchat.io:30000/rtm-api