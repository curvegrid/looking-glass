#!/bin/bash

cd "$(dirname "$0")"

exec go run . -web ../web/dist
