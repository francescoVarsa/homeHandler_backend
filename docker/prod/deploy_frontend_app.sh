#!/bin/bash

echo deploy frontend app...

FILE=./www.tar.gz
exit_code=$1

if [ -f "$FILE" ]; then
    rm -rf build
    tar -xzvf www.tar.gz

    echo restarting containers...
    docker compose restart

    echo deploy successfully finished

else 
    echo The file www.tar.gz does not exits under this path $FILE

    exit $exit_code
fi
