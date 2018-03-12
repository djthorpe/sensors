#!/bin/bash
##############################################################
# Build Linux Flavours
##############################################################

CURRENT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO=`which go`
LDFLAGS="-w -s"
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
    bme280.go
    ener314.go
    tsl2561.go
)

for COMMAND in ${COMMANDS[@]}; do
    echo "go install cmd/${COMMAND}"
    echo go install -ldflags \"${LDFLAGS}\" \"cmd/${COMMAND}\"
    go install -ldflags "${LDFLAGS}" "cmd/${COMMAND}" || exit -1
done

