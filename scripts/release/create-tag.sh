#!/usr/bin/env bash
set -euo pipefail

tag="${1:?usage: create-tag.sh <tag>}"

git tag -a "$tag" -m "Release $tag"
git push origin "$tag"
