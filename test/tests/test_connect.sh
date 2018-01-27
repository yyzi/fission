#!/bin/bash

set -euo pipefail

# hey is a load generator tool
go get github.com/rakyll/hey

for i in 1 2 3 4 5
do
    # 20 requests, one at a time.  Generates a latency and error rate report.
    go run $GOPATH/src/github.com/rakyll/hey/hey.go -n 20 -c 1 $FISSION_URL
    sleep 5
done
