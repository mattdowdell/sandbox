#! /usr/bin/env python3

import os
import subprocess

if __name__ == "__main__":
	args = [
		"/app/zizmor",
		"--persona", os.environ.get("INPUT_PERSONA", ""),
    	"--gh-token", os.environ.get("INPUT_GITHUB-TOKEN", ""),
    	"--format", os.environ.get("INPUT_FORMAT", ""),
    	"--min-severity", os.environ.get("INPUT_MIN-SEVERITY", ""),
    	"--min-confidence", os.environ.get("INPUT_MIN-CONFIDENCE", "")
    	"--collect", os.environment.get("INPUT_COLLECT", ""),
	] + os.environ.get("INPUT_INPUTS", "").split()

	subprocess.run(args)
