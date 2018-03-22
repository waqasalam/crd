#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if [ "$#" -lt 5 ] || [ "${1}" == "--help" ]; then
    cat <<EOF
Usage: $(basename $0)  <output-package> <package-path> <group> <version> <controller-path>
  <output-package> package path for generated files
  <package-path> Path to root of package 
  <group>  Group name
  <version> Version 
  <controller-path> Path where controller files are generated
EOF
    exit 0
fi

OUTPKGPATH=$1
PKGPATH=$2
GROUP=$3
VERSION=$4
CONTROLLER=$5

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
go install ${SCRIPT_ROOT}/script/

echo "Generating crds and controller"
${GOPATH}/bin/script --output-package ${OUTPKGPATH} --pkg-path ${PKGPATH} --group ${GROUP} --version ${VERSION} --controller ${CONTROLLER}  
