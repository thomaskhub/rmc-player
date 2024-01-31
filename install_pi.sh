#!/bin/bash 
#
# Installation script for installing everythign on raspbian os
# 
# How to call to install: 
# ./install_pi.sh install <wifi_password> <wifi_ssid>
# 

#
# configs
#
GIT_REPO=https://github.com/thomaskhub/rmc-player
GIT_BRANCH=main

USBMNT_REPO=https://github.com/thomaskhub/usbmnt.git
USBMNT_BRANCH=main

GO_BIN=/opt/rmc/go/bin/go
INSTALL_DIR=/opt/rmc
DEKSTOP_DIR=/usr/share/applications

GO_URL=https://go.dev/dl/go1.21.6.linux-armv6l.tar.gz
GO_NAME=go1.21.6.linux-armv6l.tar.gz

WIFI_PASSWORD=$2
WIFI_SSID=$3



#
# Functions 
#
function echoInfo() {
    echo -e "\033[0;34m[info]: $1\033[0m"
}

function echoError() {
    echo -e "\033[0;31m[error]: $1\033[0m"
}

function installPlymouth() {
    echoInfo "installing plymouth..."
    sudo apt install -y plymouth plymouth-themes

    sudo cp -r ./assets/plymouth-image /usr/share/plymouth/themes/image
    sudo plymouth-set-default-theme image 
    sudo update-initramfs -u 

    #chek if /boot/cmdline containe plymouth settings to enable splash screen, if not add them this
    # must be executed as sudo  
  
    echo "splash" | sudo tee -a /boot/cmdline.txt
    echo "plymouth.enable=1" | sudo tee -a /boot/cmdline.txt
    echo "logo.nologo" | sudo tee -a /boot/cmdline.txt
    echo "consoleblank=0" | sudo tee -a /boot/cmdline.txt
    echo "vt.global_cursor_default=0" | sudo tee -a /boot/cmdline.txt
    echo "loglevel=1" | sudo tee -a /boot/cmdline.txt
    echo "quiet" | sudo tee -a /boot/cmdline.txt
    echo "plymouth.ignore-serial-consoles" | sudo tee -a /boot/cmdline.txt
}

function installRmc() {
    echoInfo "installing rmc..."
    workDir=$(pwd)

    cd $INSTALL_DIR

    git clone $GIT_REPO -b $GIT_BRANCH
    sudo apt install -y  libmpv2 libmpv-dev libsdl2-2.0-0 xorg xserver-xorg \
        xinit x11-xserver-utils dbus pulseaudio

    # compile rmc player ![cd into the player directory]
    echoInfo "compiling rmc player..."
    cd $INSTALL_DIR/rmc-player
    CGO_ENABLED=1 GOARCH=arm64 $GO_BIN build -o rmc-player
    
    sudo cp ./assets/rmc.service /etc/systemd/system/rmc.service
    sudo systemctl enable rmc.service
    sudo systemctl start rmc.service

    sudo chmod +x /opt/rmc/rmc-player/assets/pi_run.sh 

    cd $workDir
}

function installGit() {
  #check if git is installed
    git --version &> /dev/null
    if [[ $? -eq 127 ]]; then
        echoInfo "git not installed, installing..."
        sudo apt install -y git
    fi
}

function installGo() {
  #check if go file exits or not
    workDir=$(pwd)
    cd $INSTALL_DIR
    if [[ ! -f $GO_BIN ]]; then
        wget -nc $GO_URL
        tar -xzf $GO_NAME
        rm $GO_NAME
    fi
    cd $workDir
}

function installUsbmnt() {
    echoInfo "installing usbmnt..."
    workDir=$(pwd)
    cd $INSTALL_DIR
    git clone $USBMNT_REPO -b $USBMNT_BRANCH
    cd $INSTALL_DIR/usbmnt
    CGO_ENABLED=1 GOARCH=arm64 $GO_BIN build -o usbmnt
    cd $workDir

    sudo cp /opt/rmc/usbmnt/usbmnt.service /etc/systemd/system/usbmnt.service
    sudo systemctl enable usbmnt.service
    sudo systemctl start usbmnt.service
}

function installHotspot() {
    echoInfo "installing hotspot..."
    sudo apt install -y  hostapd dnsmasq dhcpcd5 

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

    echo "$WIFI_PASSWORD"  | sudo tee  /boot/wifi_pw.txt
    echo "$WIFI_SSID"  | sudo tee /boot/wifi_ssid.txt    
}

function setupAutologin(){
    cat > /tmp/autologin.conf << EOF
[Service]
ExecStart=
ExecStart=-/sbin/agetty --autologin $USER --noclear %I \$TERM
EOF

    sudo cp /tmp/autologin.conf /etc/systemd/system/getty@tty1.service.d/autologin.conf
}

#
# Install routine 
#
function all() {
    CWD=$(pwd)
    echoInfo "create installation directories..."

    if [[ -d $INSTALL_DIR/rmc-player ]]; then
        rm -rf $INSTALL_DIR/rmc-player
    fi;

    #check if usbmount exists if so delete it to start clean 
    if [[ -d $INSTALL_DIR/usbmnt ]]; then
        rm -rf $INSTALL_DIR/usbmnt
    fi;

    sudo mkdir -p $INSTALL_DIR
    sudo chown -R $USER $INSTALL_DIR
    sudo chmod -R 777 $INSTALL_DIR

    sudo cp ./assets/cmdline.txt /boot/cmdline.txt

    # cd $INSTALL_DIR

    installGit

    installGo

    installRmc

    installUsbmnt
   
    installHotspot
       
    #setupAutolog

    installPlymouth

    echoInfo "installation complete"
}

#get the first argument of the script and call the corresponding function. 
# for all, installHotspot, installPlymouth, installRmc, installUsbmnt, installGo, installGit, setupAutologin
# the argument name is the name of the function. Check if function exists if not print error

if [[ -z $1 ]]; then
    echo "please specify an argument: installHotspot, installPlymouth, installRmc, installUsbmnt, installGo, installGit, setupAutologin, all"
    exit 1
fi

# if argument is installHotspot print hello 
if [[ $1 == "installHotspot" ]]; then
     #if password or ssid is not set exit printing a hint of what needs to be added 
    if [[ -z $WIFI_PASSWORD || -z $WIFI_SSID ]]; then
        echo "Usage: install_pi.sh install <wifi_password> <wifi_ssid>"
        exit 1
    fi
fi


# call the function now
$1 $@



