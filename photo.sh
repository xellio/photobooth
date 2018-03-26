#!/bin/bash
NAME=${1}
PHOTO_PATH=${2}
THUMB_PATH=${3}

env LANG=C gphoto2 --port usb: --capture-image-and-download --filename $PHOTO_PATH$NAME
convert -thumbnail 200 $PHOTO_PATH$NAME $THUMB_PATH$NAME