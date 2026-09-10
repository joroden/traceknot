#!/usr/bin/env bash
set -euo pipefail

latest_tag=$(git tag -l 'v*.*.*' | sort -V | tail -n1)

range=""
if [ -n "$latest_tag" ]; then
  range="${latest_tag}..HEAD"
fi

commit_hashes=$(git log $range --format=%H)

if [ -z "$commit_hashes" ]; then
  echo "No commits since ${latest_tag:-the beginning of history}; nothing to release" >&2
  exit 1
fi

bump_rank() {
  case "$1" in
    major) echo 3 ;;
    minor) echo 2 ;;
    patch) echo 1 ;;
    *) echo 0 ;;
  esac
}

bump="none"

while IFS= read -r hash; do
  subject=$(git log -1 --format=%s "$hash")
  body=$(git log -1 --format=%B "$hash")

  candidate="none"
  if echo "$subject" | grep -qE '^[a-zA-Z]+(\([^)]*\))?!:' || echo "$body" | grep -qE '^BREAKING CHANGE:'; then
    candidate="major"
  elif echo "$subject" | grep -qE '^feat(\([^)]*\))?:'; then
    candidate="minor"
  elif echo "$subject" | grep -qE '^fix(\([^)]*\))?:'; then
    candidate="patch"
  fi

  if [ "$(bump_rank "$candidate")" -gt "$(bump_rank "$bump")" ]; then
    bump="$candidate"
  fi
done <<< "$commit_hashes"

if [ "$bump" = "none" ]; then
  echo "No fix/feat commits since ${latest_tag:-the beginning of history}; nothing to release" >&2
  exit 1
fi

base_version="${latest_tag:-v0.0.0}"
IFS='.' read -r major minor patch <<< "${base_version#v}"

case "$bump" in
  major) major=$((major + 1)); minor=0; patch=0 ;;
  minor) minor=$((minor + 1)); patch=0 ;;
  patch) patch=$((patch + 1)) ;;
esac

tag="v${major}.${minor}.${patch}"

if git rev-parse "$tag" >/dev/null 2>&1; then
  echo "Computed tag $tag already exists" >&2
  exit 1
fi

echo "$tag"
