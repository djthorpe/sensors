#!/bin/bash
##############################################################
# Build sensors which communicate over the SPI bus
##############################################################

CURRENT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO=`which go`
LDFLAGS="-w -s"
TAGS="i2c"
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
    tsl2561.go
    bme280.go 
)

echo "tags=${TAGS}"
for FILES in ${COMMANDS[@]}; do
  EXEC=`dirname ${FILES}` 
  DIR=`dirname ${FILES}` 
  if [ ${EXEC} == "." ] ; then
    EXEC=`basename -s .go ${FILES}`
  fi
  SOURCES=`basename ${FILES}`
  echo "go install ${EXEC}"
  go build -ldflags "${LDFLAGS}" -o "${GOBIN}/${EXEC}" -tags ${TAGS} "${CURRENT_PATH}/${DIR}/"${SOURCES} || exit -1
done
