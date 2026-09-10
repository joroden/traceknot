#!/usr/bin/env bash
set -euo pipefail

tag="${1:?usage: create-github-release.sh <tag>}"

gh release create "$tag" --title "$tag" --generate-notes
