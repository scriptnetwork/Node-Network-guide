#!/bin/bash

# Build a docker image for a Script node.
# Usage: 
#    integration/docker/node/build.sh
#
# After the image is built, you can create a container by:
#    docker stop script_node
#    docker rm script_node
#    docker run -e SCRIPT_CONFIG_PATH=/script/integration/scriptnet/node --name script_node -it script
set -e

SCRIPTPATH=$(dirname "$0")

echo $SCRIPTPATH

if [ "$1" =  "force" ] || [[ "$(docker images -q script 2> /dev/null)" == "" ]]; then
    docker build -t script -f $SCRIPTPATH/Dockerfile .
fi


