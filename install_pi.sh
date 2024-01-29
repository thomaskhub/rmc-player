#!/bin/bash 
#
# Installation script for installing everythign on raspbian os
#

#
# configs
#
GIT_REPO=https://github.com/thomaskhub/rmc-player
GIT_BRANCH=main

GO_BIN=/opt/rmc/go/bin/go
INSTALL_DIR=/opt/rmc
DEKSTOP_DIR=/usr/share/applications

GO_URL=https://go.dev/dl/go1.21.6.linux-armv6l.tar.gz
GO_NAME=go1.21.6.linux-armv6l.tar.gz

#
# Functions 
#
function echoInfo() {
    echo -e "\033[0;34m[info]: $1\033[0m"
}

function echoError() {
    echo -e "\033[0;31m[error]: $1\033[0m"
}

function installPi () {
    echoInfo "installing pi specific things to make it run as a standaline media center"
    sudo apt install -y  libmpv2 libmpv-dev libsdl2-2.0-0 xorg xserver-xorg \
        xinit hostapd dnsmasq dhcpcd5 plymouth plymouth-themes rpd-plym-splash \
        x11-xserver-utils dbus
    
    #install configuration for enabling wifi hot spot on wlan0 with ip 192.168.90.1, and channel 11
    echoInfo "install wifi hot spot"

    sudo systemctl stop NetworkManager
    sudo systemctl disable NetworkManager

    sudo systemctl stop wpa_supplicant
    sudo systemctl disable wpa_supplicant
  
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

    #
    # setup splach screen to show logo during boot 
    #
    # disable splash screen [this is taken from raspi-config script]
    cat > /etc/systemd/system/getty@tty1.service.d/autologin.conf << EOF
[Service]
ExecStart=
ExecStart=-/sbin/agetty --autologin $USER --noclear %I \$TERM
EOF

    cp ./icons/splash.png /usr/share/plymouth/themes/pix/splash.png
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

    cd $INSTALL_DIR

    #check if git is installed
    git --version &> /dev/null
    if [[ $? -eq 127 ]]; then
        echoInfo "git not installed, installing..."
        sudo apt install -y git
    fi

    echoInfo "Cloning git repository..."
    git clone $GIT_REPO -b $GIT_BRANCH

    #check if go file exits or not
    if [[ ! -f $GO_BIN ]]; then
        echoInfo "go version $version is not ok, installing it temorarily..."
        wget -nc $GO_URL
        tar -xzf $GO_NAME
        rm $GO_NAME
    fi

    sudo apt install -y  libmpv2 libmpv-dev libsdl2-2.0-0

    # compile rmc player 
    echoInfo "compiling rmc player..."
    cd $INSTALL_DIR/rmc-player
    CGO_ENABLED=1 GOARCH=arm64 $GO_BIN build -o rmc-player
    cd $CWD

    #install pi specific things
    installPi
 
    # now we are done 

    echoInfo "installation complete"
}

install
