#!/bin/bash

#
# This is the start up script for the pi application. It performt the following 
# things: 
#
#   - host sport setup 
#   - start the RMC player application
#       - if we have two monitors connected: 
#           - 1st monitor will run the rmc-app to control the player
#           - 2nd monitor will run the rmc-player application
#       - if we have only one monitor connected: 
#           - we will display only the player, controls with application must be done 
#             from another laptop or android phone
#

# setup and start the hotspot on every boot
echo "Creating hotspot"


# first we setup the hotspot name and passowrd. byu reading it froom config files 
# in the boot directory. This way users can easily change user name password 
# by pluggin the sd card into any computer with fat32 support
WIFI_PASSWORD=$(head -n 1 /boot/wifi_pw.txt)
WIFI_SSID=$(head -n 1 /boot/wifi_ssid.txt)

cp /opt/rmc/rmc-player/assets/hostapd.conf /tmp/hostapd.conf
sed -i "s/^wpa_passphrase=.*$/wpa_passphrase=${WIFI_PASSWORD}/" /tmp/hostapd.conf
sed -i "s/^ssid=.*$/ssid=${$WIFI_SSID}/" /tmp/hostapd.conf
sudo cp /rmp/hostapd.conf /etc/hostapd/hostapd.conf


sudo systemctl stop wpa_supplicant
sudo systemctl stop NetworkManager
sudo systemctl stop dnsmasq
sudo systemctl stop hostapd

sudo ip addr flush wlan0
sudo ip link set dev wlan0 down
sudo ip addr add 192.168.90.1/24 brd + dev wlan0
sudo ip link set dev wlan0 up

sudo systemctl unmask dnsmasq
sudo systemctl unmask hostapd
sudo systemctl start dnsmasq
sudo systemctl start hostapd


# start the player and ensure it comes back up when terminated
# This is blocking so nothing should come after this loop
# first move in the player dir so that assets / configs are being picked up 
# from default location
cd /opt/rmc/rmc-player
while true; do
    startx $(pwd)/rmc-player 
    sleep 5
done