#!/usr/bin/env bash

package=$1
package_split=(${package//\// })
package_name=${package_split[-1]}


CURDIR=`/bin/pwd`
BASEDIR=$(dirname $0)
ABSPATH=$(readlink -f $0)
ABSDIR=$(dirname $ABSPATH)

PLATFORMS=""
#PLATFORMS="$PLATFORMS linux/amd64 linux/386"
PLATFORMS="$PLATFORMS linux/amd64"




type setopt >/dev/null 2>&1

SCRIPT_NAME=`basename "$0"`
FAILURES=""
CURRENT_DIRECTORY=${PWD##*/}
OUTPUT="$package_name" # if no src file given, use current dir name


for PLATFORM in $PLATFORMS; do
  echo "platform: $PLATFORM"
  GOOS=${PLATFORM%/*}
  echo "GOOS: $GOOS"
  GOARCH=${PLATFORM#*/}
  echo "GOARCH: $GOARCH"
  OUTPUT_DIR="${ABSDIR}/build/dero_${GOOS}_${GOARCH}"
  echo "OUTPUT_DIR: $OUTPUT_DIR"
  BIN_FILENAME="${OUTPUT}-${GOOS}-${GOARCH}"
  echo "BIN_FILENAME: $BIN_FILENAME"
echo  mkdir -p $OUTPUT_DIR
  #if [[ "${GOOS}" == "windows" ]]; then BIN_FILENAME="${BIN_FILENAME}.exe"; fi
  CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o $OUTPUT_DIR/${BIN_FILENAME} $package"
  echo "cmd: ${CMD}"
  eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
done

# ARM64 builds only for linux
if [[ $PLATFORMS_ARM == *"linux"* ]]; then
  GOOS="linux"
  GOARCH="arm64"
  OUTPUT_DIR="${ABSDIR}/build/dero_${GOOS}_${GOARCH}"
  CMD="GOOS=linux GOARCH=arm64 go build -o $OUTPUT_DIR/${OUTPUT}-linux-arm64 $package"
  echo "${CMD}"
  eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
fi




# eval errors
if [[ "${FAILURES}" != "" ]]; then
  echo ""
  echo "${SCRIPT_NAME} failed on: ${FAILURES}"
  exit 1
fi
