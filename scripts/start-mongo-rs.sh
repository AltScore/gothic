#!/bin/bash
# A sample to start a local mongo
docker run -it --rm --name mongo \
  -p 27017:27017 \
  -e MONGODB_ADVERTISED_HOSTNAME=localhost \
  -e ALLOW_EMPTY_PASSWORD="yes" \
  -e MONGODB_REPLICA_SET_NAME=rs0 \
  -e MONGODB_REPLICA_SET_MODE=primary \
  bitnami/mongodb:4.4
