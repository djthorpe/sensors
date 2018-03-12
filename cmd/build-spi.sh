#!/bin/bash
##############################################################
# Build sensors which communicate over the SPI bus
##############################################################

CURRENT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO=`which go`
LDFLAGS="-w -s"
TAGS="spi"
cd "${CURRENT_PATH}/.."

##############################################################
# Sanity checks

if [ ! -d ${CURRENT_PATH} ] ; then
  echo "Not found: ${CURRENT_PATH}" >&2
  exit -1
fi
if [ "${GO}" == "" ] || [ ! -x ${GO} ] ; then
  echo "go not installed or executable" >&2
  exit -1
fi

##############################################################
# Install

COMMANDS=(
    rfm69/*.go
    bme280.go 
)

for COMMAND in ${COMMANDS[@]}; do
  EXEC=`dirname ${COMMAND}`
  if [ ${EXEC} == "." ] ; then
    EXEC=`basename -s .go ${COMMAND}`
  fi
  echo "go install ${EXEC}"
  go build -ldflags "${LDFLAGS}" -o "${GOBIN}/${EXEC}" -tags "${TAGS}" "${FILES}" || exit -1
done
