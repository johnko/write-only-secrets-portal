#!/usr/bin/env bash
set -ex

open "http://localhost:8888/"
python -m http.server 8888 --bind 127.0.0.1 --directory ./src/
