# Overview

Remote Media Center (RMC) is a video player which can be controlled by
an http interface. Check out ./server.go file to see the http routes available.
The idea is to have a media player with very minimal on screen display components
and have a remote interface app running e.g. on Android to control it.

The idea is to separate the controls from the player, which is useful for Kiosk
applications etc.

The base is the mpv player (libmpv) packaged around an http interface
written in golang. Most of the code runs on all platforms, except a few
functions which have been added to be used on raspberry pi.

# Installation on Ubuntu 22 / 23

We clone the code, then build it and install it. This way it will run on many
different platforms. For different distros you might need to adjust the install_ubuntu script slighlty.

1. download or clone the player rep locally
2. execute the install script, which will instal everything under /opt/rmc

```bash
bash ./install_ubuntu.sh
```

# Installation on PI

see rmc-pi repo for installation of the player on the PI.
Basically it needs go to compile it, and libmpv and libsdl2 libraries.
It only runs on the 64b bit version of Raspbian OS

```bash
#compile go code on PI (arm)
CGO_ENABLED=1 GOARCH=arm64 $GO_BIN build -o rmc-player
```

# Installation on Windows

- Install libmpv and libsdl and compile it (you have to figure out the details for now :)
- use it as is or create executable files and share it.
- TODO: we can add the windows executable files as release files to the repo

# Installation on Mac

- This still needs to be figured out
