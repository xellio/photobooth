# Photobooth

[![go report card](https://goreportcard.com/badge/github.com/xellio/photobooth "go report card")](https://goreportcard.com/report/github.com/xellio/photobooth)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/xellio/photobooth?status.svg)](https://godoc.org/github.com/xellio/photobooth)

## Ambition
I wanted to build a simple, lightweight and easy to use photobooth application.
In my case it runns on a Rasperry Pi 3 - but it should run on any system/hardware matching the requirements.

The application consists of two pages, the index (`/`) and the gallery (`/gallery`).
The photobooth should only show the index page.
The gallery is designed for accessing the photos from other devices/people in your network.

## Note
This application is designed to work at small events in a more or less controlled environment.

Due to security, privacy and performance concerns, you should not use this software (without further checks and modifications) in uncontrolled environments.

## Requirements
- [go](https://golang.org/)
- [gphoto2](http://gphoto.org/)
- imagemagick
- convert

## Installation
Running `make` will do the work for you.

Alternatively you can build and run the application on your own
```
# 1. create the desired directories for saving the photos and thumbnails
# 2. build the application from the applications root folder
go build -ldflags="-s -w" -o photobooth ./cmd/main.go
# 3. run the application
./photobooth --port=1122 --limit=100 --photopath=/the/path/for/photos/ --thumbpath=/the/path/for/thumbnails/
# 4. point your browser to http://localhost:1122/ (instead of localhost - your photobooth's IP)
```

## Configuration
You can configure the save path in the `config` file, located in the root directory.
```
# the location for saving the photos
PHOTO_DIR = $(MAKEFILEDIR)/photos/

# the location for saving the thumbnails
THUMB_DIR = $(MAKEFILEDIR)/photos/thumbs/

# limit of loaded photos (in the frontend)
IMAGE_LIMIT = 100

# the port for the application
APPLICATION_PORT = 1122
```

## Setup
My setup requires a Raspberry Pi 3 running a headless raspbian stretch, a [gphoto2 supported camera](http://gphoto.org/doc/remote/) (for me it is a Nikon D5000) and a Google Nexus 7 for accessing and controlling the frontend from inside the booth.
