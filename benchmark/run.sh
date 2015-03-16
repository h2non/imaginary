#!/usr/bin/env bash

go build

./imgine -p 8088 > /dev/null &

boom -n 1000 -c 20 http://localhost:8088 -m POST -T "image/jpeg" -d 