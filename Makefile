lk: main.go
	go build
	xdg-open http://0.0.0.0:3000
	./lk /home/hendry/media/scans/wedding_20110327132046

clean:
	rm -f lk
