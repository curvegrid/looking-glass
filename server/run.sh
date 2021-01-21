#!/bin/bash

cd "$(dirname "$0")"

CONFIG_FILE="looking-glass.json"

for i in "$@"
do
    case $i in
        -c=*)
            CONFIG_FILE="${i#*=}"
            shift
            ;;
    esac
done

exec go run . -web ../web/dist -c $CONFIG_FILE
