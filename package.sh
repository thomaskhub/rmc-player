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
    filename="rmc.exe"
    # arch="amd64"
    # goplatform="windows"
elif [[ "$platform" == "linux" ]]; then 
    filename="rmc"
    # options="-rpath='./linux.lib'"
    cp -r ./linux.lib dist/$platform
    # arch="amd64"
    # goplatform="linux"
elif [[ "$platform" == "darwin" ]]; then 
    filename="rmc"
    # arch="amd64"
    # goplatform="darwin"
elif [[ "$platform" == "pi" ]]; then 
    filename="rmc"
    # arch="arm64"
    # goplatform="linux"
fi

# Compile the go application for the different operating systems using cross compiler option 
go mod tidy
go build -o dist/$platform/$filename . $options
# GOOS=$goplatform GOARCH=$arch go build -o dist/$platform/$filename .

# When we are on linux also create an appImage which might be better portable
# install all tools if needed 
if [[ "$platform" == "linux" ]]; then 
  wdir=$(pwd)
  # echo "going to create appimage $wdir"

  # mkdir -p ./dist/rmc.AppDir/{bin,downloads,lib} 
  # elfPath="$wdir/dist/$platform/rmc"

  # cp ./dist/$platform/rmc ./dist/rmc.AppDir/bin

  # cd ./dist/rmc.AppDir 
  # for lib in $(ldd $elfPath | awk '{print $3}'); do 
  #   dirname=$(dirname $lib)
  #   mkdir -p .$dirname
  #   cp "$lib" .$lib; 
  # done  # Copy dependencies

  # mkdir -p ./lib64
  # cp /lib64/ld-linux-x86-64.so.2 ./lib64

  

  cd $wdir


  # mkdir -p dist/appImage/tools
  # cd dist/appImage

  # #Install appimage-builder if not already availanle
  # wget https://github.com/AppImageCrafters/appimage-builder/releases/download/v1.1.0/appimage-builder-1.1.0-x86_64.AppImage

  # appimage-builder \
  #   --appimage-extract=true \
  #   --executable ../linux/rmc
  
  # appimage-builder

  cd $wdir
fi

# Create zip archives of the distributions in the dist dir 
cd dist
zip  rmc-player-$platform.zip $platform
cd ..






