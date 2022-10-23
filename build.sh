#!/bin/bash
mkdir -p ./tmp/builds/
chmod +777 ./tmp/builds/
outputName=./tmp/builds/"$2"-"$3"
go build -o "$outputName" "$1"
if [ "$?" == 1 ];
then
    echo "build failed; path: $1; output: $outputName"
    exit 1
fi
chmod +x "$outputName"
echo "$outputName"