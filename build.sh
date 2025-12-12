#!/usr/bin/env bash

# A multiarch builder is needed to run this script:
# docker buildx create --name multiarch --driver docker-container --use

set -euo pipefail

VERSIONS_DIR="versions"
IMAGE="alexycodes/resticprofile"
LATEST="0.32.0"

if [[ ! -d "$VERSIONS_DIR" ]]; then
  echo "Directory not found: $VERSIONS_DIR" >&2
  exit 1
fi

start=$(date +%s)

for dir in $(printf '%s\n' "$VERSIONS_DIR"/*/ | sort -V); do
  [[ -d "$dir" ]] || continue

  version="$(basename "$dir")"
  dockerfile="${dir}Dockerfile"
  platform_file="${dir}platform"

  if [[ ! -f "$dockerfile" ]]; then
    echo "No Dockerfile in $dir" >&2
    exit 1
  fi

  if [[ ! -f "$platform_file" ]]; then
    echo "No platform file in $dir" >&2
    exit 1
  fi

  platform=$(<"$platform_file")

  if [[ "${version}" == "${LATEST}" ]]; then
    echo -e "\nBuilding $IMAGE:$version (latest) for platform(s) $platform\n"

    docker buildx build \
      --no-cache \
      --platform "$platform" \
      -t "$IMAGE:$version" \
      -t "$IMAGE:latest" \
      -f "$dockerfile" \
      --push \
      .
  else
    echo -e "\nBuilding $IMAGE:$version for platform(s) $platform\n"

    docker buildx build \
      --no-cache \
      --platform "$platform" \
      -t "$IMAGE:$version" \
      -f "$dockerfile" \
      --push \
      .
  fi
done

end=$(date +%s)
runtime=$(( end - start ))

printf 'Finished in %02d:%02d:%02d\n' \
       $(( runtime/3600 )) $(( (runtime/60)%60 )) $(( runtime%60 ))
