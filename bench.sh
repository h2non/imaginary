#!/bin/bash
#
# Simple benchmark test suite
#
# You must have installed vegeta:
# go get github.com/tsenart/vegeta
#

host="http://localhost:8088"

./imaginary -p 8088 & > /dev/null
pid=$!

echo "Enlarge"
echo "POST $host/enlarge?width=1200&height=800" | vegeta attack \
  -duration=10s \
  -rate=15 \
  -body="./fixtures/large.jpg" \ | vegeta report

echo
sleep 1

echo "Resize --------------------------------------"
echo "POST $host/resize?width=200" | vegeta attack \
  -duration=10s \
  -rate=20 \
  -body="./fixtures/large.jpg" \ | vegeta report

echo 
sleep 1

echo "Extract -------------------------------------"
echo "POST $host/extract?top=50&left=50&areawidth=200&areaheight=200" | vegeta attack \
  -duration=10s \
  -rate=20 \
  -body="./fixtures/large.jpg" \ | vegeta report

echo
sleep 1

echo "Rotate -------------------------------------"
echo "POST $host/rotate?rotate=180" | vegeta attack \
  -duration=10s \
  -rate=20 \
  -body="./fixtures/large.jpg" \ | vegeta report

kill -9 $pid
