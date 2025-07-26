#! /usr/bin/env bash

set -eu -o pipefail

if [ $# -ne 1 ]; then
  echo "Usage: $0 <TOOL>"
  exit 1
fi

TOOL=$1
FILENAME=.bingo/$1.mod

echo "version=`grep --no-filename require $FILENAME | head -n1 | awk '{print $3}'`"
