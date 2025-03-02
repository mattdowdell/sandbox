#!/bin/bash

cmd=(/app/zizmor $@ ${INPUT_INPUTS[@]})
echo "Running command: ${cmd[@]}"
${cmd[@]}
