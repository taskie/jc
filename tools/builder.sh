#!/bin/bash
set -eu
set -o pipefail

CMD=
INSTALL=
while getopts 'ci' OPTS
do
    case "$OPTS" in
        c)  # Cmd
            CMD=1 ;;
        i)  # Install
            INSTALL=1 ;;
    esac
done
shift $((OPTIND - 1))

VERSION="$(git describe --tags --abbrev=0 || :)"
REVISION="$(git rev-parse --short HEAD || :)"
LDFLAGS="-X 'jc.version=${VERSION}' -X 'jc.revision=${REVISION}'"

if (( INSTALL )); then
    go install -ldflags "$LDFLAGS"
else
    go build -ldflags "$LDFLAGS"
fi
