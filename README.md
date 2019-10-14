[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![CircleCI](https://circleci.com/gh/danielpaulus/quicktime_video_hack.svg?style=svg)](https://circleci.com/gh/danielpaulus/quicktime_video_hack)
[![codecov](https://codecov.io/gh/danielpaulus/quicktime_video_hack/branch/master/graph/badge.svg)](https://codecov.io/gh/danielpaulus/quicktime_video_hack)
[![Go Report](https://goreportcard.com/badge/github.com/danielpaulus/quicktime_video_hack)](https://goreportcard.com/report/github.com/danielpaulus/quicktime_video_hack)

###  Operating System indepedent implementation for Quicktime Screensharing for iOS devices :-)
This repository contains all the code you will need to grab and record video and audio from your iPhone or iPad 
without needing one of these expensive MacOS X computers :-D
It probably does something similar to what `QuickTime` and `com.apple.cmio.iOSScreenCaptureAssistant` are doing on MacOS.

Progress:
1. [~~Write a prototype decoder in Java to reconstruct a Video from a USB Dump taken with Wireshark~~](https://github.com/danielpaulus/quicktime_video_hack/tree/JavaCMSampleBufDecoder/java-x264-decoder) 
2. ~~Create a Golang App to successfully enable and grab the first video stream data from the iOS Device on Linux~~
3. ~~port decoder from Java to Golang~~
3. Make the `go run main.go dumpraw` work on the first execution (currently you have to run it twice)
4. ~~Add support for CMSync.h style CMClock packets~~ 
5. ~~Get a continuous Feed from the Device~~
6. Make the dumpraw command work without having to unplug the device every time
7. ~~Parse Async Feed packets correctly~~ 
8. ~~Generate SPS and PPS h264 NaLus~~ 
9. ~~Fix last issue to get a working stream of FEEDs, currently my NEED packets are rejected~~
10. Create a .h264 file drom the dumpraw command 
11. Generate GStreamer compatible x264 stream probably by wrapping the NaLus in RTP headers
12. Complete packet documentation

Extra Goals:
1. [Port to Windows](https://github.com/danielpaulus/quicktime_video_hack/tree/windows/windows) (I don't know why, but still people use Windows nowadays)



### MAC OS X LIBUSB -- IMPORTANT
1. What works:
 You can enable the QuickTime config and discover QT capable devices with `qvh devices` and  `qvh activate` 

2. What does not work
 `qvh dumpraw` won't work on MAC OS because the binary needs to be codesigned with `com.apple.ibridge.control`
 apparently that is a protected Entitlement that I have no idea how to use or sign my binary with. 

2. Make sure to use either this fork `https://github.com/GroundControl-Solutions/libusb`
   or a LibUsb version BELOW 1.0.20 or iOS devices won't be found on Mac OS X.
   [See Github Issue](https://github.com/libusb/libusb/issues/290)

### How does the Mac get live video from iOS?
Basically the communication between Mac OSX and the iOS device is really just a bunch of serialized 
[CoreMedia-Framework](https://developer.apple.com/documentation/coremedia) classes. 
First there is a negotiation of CMFormatDescription objects and creation of CMClocks for time sync.
After that the host periodically sends `need` packets. And sometimes receives CMSyncs to keep the Clocks synchronized.
The device will respond with a series of CMSampleBuffers containing
raw h264 [Nalus](https://en.wikipedia.org/wiki/Network_Abstraction_Layer) 
Those can be forwarded to ffmpeg or gstreamer easily :-D
