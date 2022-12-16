#!/bin/bash
mkdir -p "$4"
chmod +777 "$4"
outputName="$4"/"$2"-"$3"
cd "$1" || echo "build failed could not change directory"
go build -o "$outputName"
if [ "$?" == 1 ];
then
    echo "build failed; path: $1; output: $outputName"
    exit 1
fi
chmod +x "$outputName"
echo "$outputName"