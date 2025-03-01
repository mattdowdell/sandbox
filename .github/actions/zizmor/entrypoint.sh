#!/bin/bash

set -euo pipefail
set -x

export GH_TOKEN="${INPUT_GITHUB_TOKEN}"

env | sort
whoami

/app/zizmor \
	--persona "${INPUT_PERSONA:-regular}" \
	--format "${INPUT_FORMAT:-plain}" \
	--min-severity "${INPUT_MIN_SEVERITY:-unknown}" \
	--min-confidence "${INPUT_MIN_CONFIDENCE:-unknown}" \
	--collect "${INPUT_COLLECT:-all}" \
	${INPUT_INPUTS[@]}
