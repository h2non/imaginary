#!/bin/bash
#
# Simple benchmark test suite
#
# You must have installed vegeta:
# go get github.com/tsenart/vegeta
#

# Default port to listen
port=8088

# Start the server
./bin/imaginary -p $port & > /dev/null
pid=$!

suite() {
  echo "$1 --------------------------------------"
  echo "POST http://localhost:$port/$2" | vegeta attack \
    -duration=30s \
    -rate=50 \
    -body="./testdata/large.jpg" \ | vegeta report
  sleep 1
}

# Run suites
suite "Crop" "crop?width=800&height=600"
suite "Resize" "resize?width=200"
#suite "Rotate" "rotate?rotate=180"
#suite "Enlarge" "enlarge?width=1600&height=1200"
suite "Extract" "extract?top=50&left=50&areawidth=200&areaheight=200"

# Kill the server
kill -9 $pid
