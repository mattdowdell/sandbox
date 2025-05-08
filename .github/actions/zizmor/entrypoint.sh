#!/bin/sh

cmd="/usr/bin/zizmor $@ ${INPUT_INPUTS}"
echo "Running command: ${cmd}"
${cmd}
