#!/bin/bash
set -eo pipefail

org_name=$1
push=0
if [[ $2 == "--push" ]]; then
  push=1
fi

tags=(latest)
if [[ -n $GITHUB_REF_NAME ]]; then
  # check if ref name is a version number
  if [[ $GITHUB_REF_NAME =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    major_version=$(echo "$GITHUB_REF_NAME" | cut -d. -f1)
    minor_version=$(echo "$GITHUB_REF_NAME" | cut -d. -f1,2)
    tags+=("$major_version" "$minor_version")
  fi
  sanitized="$(echo $GITHUB_REF_NAME | sed 's/[^a-zA-Z0-9.-]\+/-/g')"
  tags+=("$sanitized")
fi
echo "Tags: ${tags[*]}"

dir=./build/
if [ ! -f "$dir/Dockerfile" ]; then
  echo "No Dockerfile found"
  exit 1
fi
if [ ! -f "$dir/config.sh" ]; then
  echo "No config.sh found for Dockerfile"
  exit 1
fi
source "$dir/config.sh"
if [[ -n "$org_name" ]]; then
  DOCKER_ORG="$org_name"
fi
DOCKER_REPOSITORY=$DOCKER_REGISTRY/$DOCKER_ORG/$DOCKER_IMAGE
echo "Repo: $DOCKER_REPOSITORY"
echo "Base dir: $DOCKER_BASE_DIR"
args=""
for tag in "${tags[@]}"; do
  args+=" -t $DOCKER_REPOSITORY:$tag"
done
if [[ $push -eq 1 ]]; then
  args+=" --push"
fi

docker buildx build \
  $args \
  --platform linux/386,linux/amd64,linux/arm64/v8,linux/arm/v7 \
  --provenance=false \
  -f "$dir/Dockerfile" "$DOCKER_BASE_DIR"
