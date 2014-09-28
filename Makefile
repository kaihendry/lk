lk: main.go
	go build
	xdg-open http://0.0.0.0:3000
##	./lk /home/hendry/media/scans/wedding_20110327132046
	./lk /home/hendry/media/scans/sodwana1990_20110326115618

clean:
	rm -f lk
