#!/bin/bash

docker run -it \
  --rm \
  -p 80:80 \
  -v ./bolt.db:/data/bolt.db \
  bolt:latest
