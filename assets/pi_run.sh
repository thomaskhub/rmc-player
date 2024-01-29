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
sudo systemctl stop wpa_supplicant
sudo systemctl stop NetworkManager
sudo systemctl stop dnsmasq
sudo systemctl stop hostapd

sudo ip addr flush wlan0
sudo ip link set dev wlan0 down
sudo ip addr add 192.168.90.1/24 brd + dev wlan0
sudo ip link set dev wlan0 up
  
sudo systemctl start dnsmasq
sudo systemctl start hostapd

#
# See if we need to start the controller also
#
#TODO: we should start a small python window manager 
# once xorg and the player is started it should start the 
# controller and move it to secondary display. If we do not 
# have a 2nd display it will not start the controller. 
# on the first version controller will need to run on 
# a secondary device like laptop or android phone



# start the player and ensure it comes back up when terminated
# This is blocking so nothing should come after this loop
# first move in the player dir so that assets / configs are being picked up 
# from default location
cd /opt/rmc/rmc-player
while true; do
    startx $(pwd)/rmc-player 
    sleep 5
done