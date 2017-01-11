#!/bin/bash

DIR=~/go/src/github.com/wantonsolutions/416TA/littlesys
cd $DIR

LOCALHOST=localhost
BASEPORT=19000

PEERS=3
for (( i=0; i<PEERS; i++ ))
do
    let "PORT = BASEPORT + i"
    go run littlesys.go -ipPort=$LOCALHOST:$PORT &
    echo launching node $i on port $PORT with pid $!
done
