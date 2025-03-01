#!/bin/bash

set -euo pipefail

export GH_TOKEN="${INPUT_GITHUB_TOKEN}"

/app/zizmor \
	--persona "${INPUT_PERSONA:-regular}" \
	--format "${INPUT_FORMAT:-plain}" \
	--min-severity "${INPUT_MIN_SEVERITY:-unknown}" \
	--min-confidence "${INPUT_MIN_CONFIDENCE:-unknown}" \
	--collect "${INPUT_COLLECT:-all}"
	"${INPUT_INPUTS}"
