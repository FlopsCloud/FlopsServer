#!/bin/bash

# Check if tables parameter is provided
if [ -z "$1" ]; then
    echo "Please provide tables parameter"
    echo "Usage: ./gen_model.sh <tables>"
    exit 1
fi

# Run goctl command with the provided tables
goctl model mysql datasource \
    -url=" : @tcp( . :3306)/flops_cloud" \
    -table="$1" \
    -dir ./model 
