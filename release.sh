
#!/bin/bash

VERSION=$(git describe --tags --abbrev=0)
VERSION_REGEX="^v(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)(\\-[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?(\\+[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?$"
GITHUB_TOKEN=$(cat "${HOME}"/.githubtoken )

if [[ ! ${VERSION} =~ $VERSION_REGEX ]]; then
    echo "Latest tag $VERSION must be a SemVer, exiting ..."
    exit 1
fi

if [[ -z ${GITHUB_TOKEN} ]]; then
  echo "Please set the \$GITHUB_TOKEN, exiting ..."
  exit 1
fi

docker run -it --rm \
  -v "$PWD":/src \
  -w /src \
  -e GITHUB_TOKEN="${GITHUB_TOKEN}" \
  goreleaser/goreleaser:v0.137.0-cgo release --rm-dist