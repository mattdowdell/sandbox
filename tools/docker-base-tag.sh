#! /usr/bin/env bash

# Check if the correct number of arguments are passed
if [ $# -lt 2 ]; then
  echo "Usage: $0 <Dockerfile> <image_name_1> <image_name_2> ... <image_name_n>"
  exit 1
fi

# File path to the Dockerfile
DOCKERFILE=$1

# Shift the first argument to get the image names
shift

# Check if the Dockerfile exists
if [ ! -f "$DOCKERFILE" ]; then
  echo "Error: $DOCKERFILE does not exist."
  exit 1
fi

# Loop over the image names provided as arguments
for image in "$@"; do
  # Extract the tag of the specific base image from the Dockerfile
  tag=$(grep -i "^FROM\s\+$image" $DOCKERFILE | \
    cut -d: -f2 | \
    cut -d@ -f1)

  if [ -z "$tag" ]; then
    echo "Image: $image - Tag: Not found" 1>&2
    exit 1
  else
    echo "tag=$tag"
  fi
done
