#!/bin/bash
mkdir -p ./tmp/builds/
outputName=./tmp/builds/"$2"-"$3"
go build -o "$outputName" "$1"
chmod +x "$outputName"
echo "$outputName"