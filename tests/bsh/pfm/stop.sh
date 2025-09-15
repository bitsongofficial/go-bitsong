#!/bin/bash

# Stop the chain
# Find and kill the process bitsongd

BINARY=$1

if pgrep -x $BINARY > /dev/null
then
    pkill -9 $BINARY
    echo "Stopped $BINARY"
else
    echo "$BINARY is not running"
fi