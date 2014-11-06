lk: main.go thumb.go
	go build

dist:
	gox

test: lk
	xdg-open http://0.0.0.0:3000
	./lk

clean:
	rm -rf lk ~/.cache/lk
