#!/bin/bash

#Assuming that the script is being called in respective directory
ORGANIZATION_ROOT_PATH="/Users/nmarler/go/src/TapMars"
PROTO_PKG_SOURCE="./productManager/pkg/proto"
PROTO_PKG_DESTINATION="./productManager_proxy/pkg/"

cd $ORGANIZATION_ROOT_PATH

## rsync ensures an identical replace/copy and can be used over ssh
## -a is 'archived mode', which copies faithfully from source to destination
## --delete removes extra file not in source from destination, as well as ensuring destination is identical to source
## -vh for verbose and human-readability information
rsync -a --delete -vh $PROTO_PKG_SOURCE $PROTO_PKG_DESTINATION