#!/usr/bin/env bash

go build

./imgine -p 8088 > /dev/null &
sleep 1

echo "Running resize tests"

ab -c 1 -n 1 -v 4 \
  -p benchmark/data.txt \
  -T "multipart/form-data; boundary=1234567890" \
  http://localhost:8088/resize

echo
echo "Running crop tests"
