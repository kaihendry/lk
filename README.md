# lk - local kuvat (pictures)

Simple Web Image Viewer, ideally for a LAN, in order to avoid the IOS Photos &
iCloud for sharing amongst friends and family

<img src=http://s.natalian.org/2014-11-04/1415116363_1364x748.png alt="Google Chrome 40 on a 1366x768 X220 display">
<img src=http://s.natalian.org/2014-11-04/lk-landscape.png alt="IOS Safari in landscape on an iPhone6">
<img src=http://s.natalian.org/2014-11-04/lk-portrait.png alt="IOS Safari in portrait on an iPhone6">

* <http://youtu.be/BQHzfpIEmwk>
* [Video of the author presenting lk at a Golang meetup](http://youtu.be/IIuDygqCOJE)

# Install

	go get -u github.com/kaihendry/lk

# Docker

* [lk.service](lk.service) for CoreOS's systemd to keep it going

How to run from Docker:

	docker run -it -p 3000:3000 --rm -v /YOUR/JPEG/IMAGES/:/srv/ lk

# Other local Web image viewers

* <https://github.com/songgao/gallery>
* <https://github.com/3d0c/imagio>
