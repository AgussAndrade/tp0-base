#!/bin/bash

ECHO_MESSAGE="Echo message"
SERVER_PORT=$(grep "SERVER_PORT" server/config.ini | awk -F '=' '{print $2}' | tr -d ' ')

RESPONSE=$(echo $ECHO_MESSAGE | docker run --rm --network=tp0_testing_net -i subfuzion/netcat -w 2 server $SERVER_PORT)


if [ "$RESPONSE" == "$ECHO_MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
