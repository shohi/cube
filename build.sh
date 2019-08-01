#!/usr/bin/env bash -x

# dockerize builds docker images, the first arg is tag.
function dockerize() {
  local GOVER=${GO_VERSION:-1.12.7}
  local VERSION=${1:-0.1.0}

  docker build \
    -t shohik/cube:${VERSION} \
    --build-arg GO_VERSION=${GOVER} \
    .
}

################################################################################
#####                          main entry                                  #####
################################################################################
function main() {
  local start=`date +%s`

  case "$1" in
   "docker")
     shift
     dockerize $@
     ;;
   *)
     echo "Unknown command"
  esac

  local end=`date +%s`
  local runtime=$((end-start))
  echo "processing time: $runtime second"
}

main $@

