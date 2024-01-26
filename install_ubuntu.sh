#!/bin/bash 
#
# Linux installation script to install rmc on linux (only tested on ubuntu 22 / 23)
# (for windows use the provided windows installer)
# (for mac use the provided macos dmg)
#

#
# configs
#
GIT_REPO=https://github.com/thomaskhub/rmc-player
GIT_BRANCH=main
GO_URL=https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
GO_BIN=/opt/rmc/go/bin/go
INSTALL_DIR=/opt/rmc
DEKSTOP_DIR=/usr/share/applications
DESKTOP_FILE=rmc-player.desktop



#
# Functions 
#
function echoInfo() {
    echo -e "\033[0;34m[info]: $1\033[0m"
}

function echoError() {
    echo -e "\033[0;31m[error]: $1\033[0m"
}



#
# Install routine 
#

function install() {
    CWD=$(pwd)

    # create the directories for installation if they do not exist
    # /opt | /opt/go | /opt/rmc-player
    echoInfo "create installation directories..."
    rm -rf $INSTALL_DIR/rmc-player
    mkdir -p $INSTALL_DIR/{go}
    sudo chown -R $USER $INSTALL_DIR
    sudo chmod -R 777 $INSTALL_DIR

    #check if git is installed
    git --version &> /dev/null
    if [[ $? -eq 127 ]]; then
        echoInfo "git not installed, installing..."
        sudo apt install -y git
    fi

    echoInfo "Cloning git repository..."
    git clone $GIT_REPO -b $GIT_BRANCH

    #get the go version from the go installation if any
    version=$($GO_BIN version | awk '{print $3}')
    IFS='.' read -r -a array <<< "$version"

    if [[ ${array[0]} -eq 1 && ${array[1]} -ge 21 ]]; then
        echoInfo "go version $version is ok"
    else
        echoInfo "go version $version is not ok, installing it temorarily..."
        wget $GO_URL
        tar -xzf go1.21.6.linux-amd64.tar.gz
    fi

    # check if we are running on ubuntu 22 or ubuntu 23
    osVersion=$(lsb_release -rs)

    #split the version with separator .
    IFS='.' read -r -a array <<< "$version"

    #if array[0] is 22 install dependencies for ubuntu 22 if its 23 install 
    # dependencies for ubuntu 23
    if [[ ${array[0]} -eq 22 ]]; then
       sudo apt install -y  libmpv1 libmpv-dev libsdl2-2.0-0 
    elif [[ ${array[0]} -eq 23 ]]; then
        sudo apt install -y  libmpv2 libmpv-dev libsdl2-2.0-0
    else
        echoError "ubuntu version not supported"
        exit 1
    fi

    # compile rmc player 
    echoInfo "compiling rmc player..."
    cd $INSTALL_DIR/rmc-player
    $GO_BIN build -o rmc-player

    # Install Desktop file so that rmc can be started from gui
    echoInfo "installing desktop file..."
    cp ./$DESKTOP_FILE $DEKSTOP_DIR

    # now we are done 
    cd $CWD

    echoInfo "installation complete"
}


function uninstall() {
    echoInfo("uninstalling rmc...")
    #do not remove the /opt/rmc dir because we are using this for the remote install also
    rm -rf $INSTALL_DIR/rmc-player
    rm -rf $INSTALL_DIR/go

    echoInfo("removing desktop file...")
    rm -rf $DEKSTOP_DIR/$DESKTOP_FILE
}


#when first argument is install call the install function, if its uninstall call the uninstall functio
if [[ "$1" == "install" ]]; then
    install
elif [[ "$1" == "uninstall" ]]; then
    uninstall
else
    echoError("invalid argument")
    exit 1
fi




