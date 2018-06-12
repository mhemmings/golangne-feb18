#!/bin/bash -e

apex init
rm functions/hello/index.js
cp ../demo1/main.go functions/hello/main.go
