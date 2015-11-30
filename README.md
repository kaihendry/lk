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
	go install github.com/kaihendry/lk

# Docker + CoreOS / Google Compute Engine

* https://registry.hub.docker.com/u/hendry/lk/builds_history/79157/
* https://blog.golang.org/docker for Google Compute Engine information, which doesn't work for me <http://r2d2.webconverger.org/2014-11-09/gce.mp4>
* [lk.service](lk.service) for CoreOS's systemd to keep it going

# Other local Web image viewers

* <https://github.com/songgao/gallery>
* <https://github.com/3d0c/imagio>
