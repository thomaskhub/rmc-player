#!/bin/bash 



# Check if the first argument is defined
if [ -z "$1" ]; then
  echo "Error: No platform specified. Supported platforms are: linux, windows, darwin."
  exit 1
fi


# Check if the platform is valid
platform="$1"
if [[ "$platform" != "linux" && "$platform" != "windows" && "$platform" != "darwin" && "$platform" != "pi"  ]]; then
  echo "Error: Invalid platform. Supported platforms are: linux, windows, darwin, pi."
  exit 1
fi

# Use the valid platform value for further actions
echo "Using platform: $platform"


# Create the fresh distribution directory 
rm -rf dist 
mkdir -p dist/$platform
cp -r ./assets dist/$platform
cp ./config.json dist/$platform
cp ./input.conf dist/$platform

# Setup the output filename OS specific

if [[ "$platform" == "windows" ]]; then
    filename="rmc-player.exe"
    # delete config file for windows its not needed 
    rm -rc dist/$platform/config.json

elif [[ "$platform" == "linux" ]]; then 
    filename="rmc-player"
    # options="-rpath='./linux.lib'"
    cp -r ./linux.lib dist/$platform
    # arch="amd64"
    # goplatform="linux"
elif [[ "$platform" == "darwin" ]]; then 
    filename="rmc-player"
    # arch="amd64"
    # goplatform="darwin"
elif [[ "$platform" == "pi" ]]; then 
    filename="rmc-player"
    # arch="arm64"
    # goplatform="linux"
fi

# Compile the go application for the different operating systems using cross compiler option 
go mod tidy
go build -o dist/$platform/$filename . $options
# GOOS=$goplatform GOARCH=$arch go build -o dist/$platform/$filename .



# Create zip archives of the distributions in the dist dir 
# TODO: for windows we actually create msi, for linux its direct install so package is 
# not needed. For mac we have to see if we can do it similar to linux or otherwise 
#w we package it
# cd dist
# zip  rmc-player-$platform.zip $platform
# cd ..






