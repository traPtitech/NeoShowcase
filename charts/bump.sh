#!/usr/bin/env bash

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <chart name> <patch|minor|major>"
  exit 1
fi

CHART=$1
STRATEGY=$2

next_version() {
  local VERSION=$1
  # Get number parts
  local MAJOR="${VERSION%%.*}"; local VERSION="${VERSION#*.}"
  local MINOR="${VERSION%%.*}"; local VERSION="${VERSION#*.}"
  local PATCH="${VERSION%%.*}"; local VERSION="${VERSION#*.}"

  # Increase version
  if [ "$STRATEGY" = "patch" ]; then
    local PATCH=$((PATCH+1))
  elif [ "$STRATEGY" = "minor" ]; then
    local MINOR=$((MINOR+1))
    local PATCH=0
  elif [ "$STRATEGY" = "major" ]; then
    local MAJOR=$((MAJOR+1))
    local MINOR=0
    local PATCH=0
  else
    echo "Unknown strategy: $STRATEGY"
    exit 1
  fi

  NEW_RAW="$MAJOR.$MINOR.$PATCH"
}

# Get current appVersion / version from Chart.yaml
APP_VERSION=$(grep '^appVersion:' "./${CHART}/Chart.yaml" | cut -d ' ' -f 2 | tr -d '"')
VERSION=$(grep '^version:' "./${CHART}/Chart.yaml" | cut -d ' ' -f 2 | tr -d '"')

# Calculate next versions
# Set appVersion to highest version tag
NEXT_APP_VERSION=$(git describe --abbrev=0 --tags --match 'v*.*.*')
NEXT_APP_VERSION="${NEXT_APP_VERSION#v}"
echo "appVersion: $APP_VERSION => $NEXT_APP_VERSION"

next_version "$VERSION"
NEXT_VERSION="$NEW_RAW"
echo "version: $VERSION => $NEXT_VERSION"

# Update Chart.yaml
sed -i "s/^appVersion: .*$/appVersion: \"${NEXT_APP_VERSION}\"/" "./${CHART}/Chart.yaml"
sed -i "s/^version: .*$/version: \"${NEXT_VERSION}\"/" "./${CHART}/Chart.yaml"
