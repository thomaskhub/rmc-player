#!/bin/bash 
#
# Linux installation script to install rmc on linux (only tested on ubuntu 22 / 23)
# (for windows use the provided windows installer)
# (for mac use the provided macos dmg)
#


GO_BIN=/opt/rmc/go/bin/go
INSTALL_DIR=/opt/rmc
DEKSTOP_DIR=/usr/share/applications
DESKTOP_FILE=rmc-player.desktop
GO_URL=https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
GO_NAME=go1.21.6.linux-amd64.tar.gz


#
# Functions 
#
function echoInfo() {
    echo -e "\033[0;34m[info]: $1\033[0m"
}

function echoError() {
    echo -e "\033[0;31m[error]: $1\033[0m"
}

function insallPi () {
    echoInfo "installing pi specific things to make it run as a standaline media center"
    #TODO: exit is only there to ensure we do not run this if we are not on the pi until its 
    #tested properly
    exit 0

    #Xorg / X11 setup
    echoInfo "install xorg and desktop dependencies"
    sudo apt install -y xorg xserver-xorg xinit

    #install configuration for enabling wifi hot spot on wlan0 with ip 192.168.90.1, and channel 11
    echoInfo "install wifi hot spot"
  
    sudo cp ./assets/autohotspot.service /etc/systemd/system/autohotspot.service
    sudo systemctl enable autohotspot.service

    sudo cp ./assets/autohotspot.sh /usr/bin/autohotspot
    sudo chmod 755 /usr/bin/autohotspot

    #copy dnamasq.conf to the dns default dir 
    sudo cp ./assets/dnsmasq.conf /etc/dnsmasq.conf
    sudo chmod 755 /etc/dnsmasq.conf

    # copy the hostapd file to the default dir 
    sudo cp ./assets/hostapd.conf /etc/hostapd/hostapd.conf
    sudo chmod 755 /etc/hostapd/hostapd.conf

    # write a service with the name rmc.service it should run /opt/rmc/rmc-player/rmc-player -c /opt/rmc/rmc-player/config.json
    # using xinit when the system boots
    sudo cp ./assets/rmc.service /etc/systemd/system/rmc.service
    sudo systemctl enable rmc.service
    sudo systemctl start rmc.service
}


#
# Install routine 
#

function install() {
    CWD=$(pwd)

    # create the directories for installation if they do not exist
    # /opt | /opt/go | /opt/rmc-player
    echoInfo "create installation directories..."
    #check if install_dir/rmc-player exists if so remove it
    if [[ -d $INSTALL_DIR/rmc-player ]]; then
        rm -rf $INSTALL_DIR/rmc-player
    fi;

    sudo mkdir -p $INSTALL_DIR
    sudo chown -R $USER $INSTALL_DIR
    sudo chmod -R 777 $INSTALL_DIR

    mkdir -p $INSTALL_DIR/rmc-player
    
    # check if the GO_BIN exists if not get GO_URL, extract it to INSTALL_DIR/go and remove the tar file
    if [[ ! -f $GO_BIN ]]; then
        echoInfo "go not installed, installing..."
        wget -P /tmp -nc $GO_URL
        tar -xzf /tmp/$GO_NAME -C $INSTALL_DIR
        rm /tmp/$GO_NAME
    fi

    # cd $INSTALL_DIR

    # check if we are running on ubuntu 22 or ubuntu 23
    osVersion=$(lsb_release -rs)

    #split the version with separator .
    IFS='.' read -r -a array <<< "$osVersion"

    #Currently supported Ubuntu  22 | 23 but any other distro should be 
    # supported using the correct packag manger command to install SDL2 and libmpv
    # but other distros would need to be tested
    if [[ ${array[0]} -eq 20 ]]; then
    #    sudo apt install -y  libmpv1 libmpv-dev libsdl2-2.0-0
       echoError "Ubuntu 20 is not supported please upgrade to 22 / 23" 
    elif [[ ${array[0]} -eq 22 ]]; then
       sudo apt install -y  libmpv1 libmpv-dev libsdl2-2.0-0 
    elif [[ ${array[0]} -eq 23 ]]; then
        sudo apt install -y  libmpv2 libmpv-dev libsdl2-2.0-0
    else
        echoError "ubuntu version not supported"
        exit 1
    fi

    # compile rmc player 
    echoInfo "compiling rmc player..."
    $GO_BIN build -o $INSTALL_DIR/rmc-player
    
    echoInfo "copy the assets and config files"
    cp -r ./assets $INSTALL_DIR/rmc-player
    cp ./input.con $INSTALL_DIR/rmc-player
    cp ./config_ubuntu.json $INSTALL_DIR/rmc-player/config.json


    # Install Desktop file so that rmc can be started from gui if its on a pc
    echoInfo "installing desktop file..."
    sudo cp ./$DESKTOP_FILE $DEKSTOP_DIR
 
    # now we are done 
    cd $CWD

    echoInfo "installation complete"
}


function uninstall() {
    echoInfo "uninstalling rmc..."
    #do not remove the /opt/rmc dir because we are using this for the remote install also
    rm -rf $INSTALL_DIR/rmc-player
    

    echoInfo "removing desktop file..."
    rm -rf $DEKSTOP_DIR/$DESKTOP_FILE
}


#when first argument is install call the install function, if its uninstall call the uninstall functio
if [[ "$1" == "install" ]]; then
    install
elif [[ "$1" == "uninstall" ]]; then
    uninstall
else
    echo "please specify an argument: install, uninstall"
    exit 1
fi




