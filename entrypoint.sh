#!/bin/sh
set -e

if [ -f /iproutes.sh ]; then
    /iproutes.sh
fi


exec /yaml-merger