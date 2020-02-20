#!/usr/bin/env sh

set -x

go build -o commandd cmd/commandd/main.go
go build -o e2etest cmd/e2etest/main.go

./commandd -addr=":8081" echo -n foo bar baz &
pid=$!

# Give commandd some time to start up and get ready to serve requests.
sleep 3

./e2etest -count=50 -response="foo bar baz" -url="http://localhost:8081/run"
code=$?

kill ${pid}
rm -f commandd e2etest

exit ${code}
