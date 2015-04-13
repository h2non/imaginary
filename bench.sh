#!/bin/bash
#
# Simple benchmark test suite
#
# You must have installed vegeta:
# go get github.com/tsenart/vegeta
#

host="http://localhost:8088"

echo "Resize -----------------------------------"
echo "POST $host/resize?width=300" | vegeta attack \
  -duration=60s \
  -rate=10 \
  -body="./fixtures/large.jpg" \ | vegeta report

exit 0
echo 

echo "Extract -------------------------------------"
echo "POST $host/extract?top=50&left=50&areawidth=200&areaheight=200" | vegeta attack \
  -duration=10s \
  -rate=20 \
  -body="./fixtures/large.jpg" \ | vegeta report

echo

echo "Rotate -------------------------------------"
echo "POST $host/rotate?rotate=180" | vegeta attack \
  -duration=10s \
  -rate=20 \
  -body="./fixtures/large.jpg" \ | vegeta report