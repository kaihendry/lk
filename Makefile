lk: main.go thumb.go
	go build

dist:
	gox

clean:
	rm -rf lk ~/.cache/lk
