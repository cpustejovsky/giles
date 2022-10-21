#!/bin/bash
mkdir -p ./tmp/builds/
chmod +x ./tmp/builds/
outputName=./tmp/builds/"$2"-"$3"
go build -o "$outputName" "$1"
chmod +x "$outputName"
echo "$outputName"