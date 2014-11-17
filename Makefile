gitVersion := $(shell git rev-parse --short HEAD)

lk: main.go thumb.go
	go build -ldflags "-X main.gitVersion $(gitVersion)"

dist:
	gox -ldflags="-X main.gitVersion $(gitVersion)"

clean:
	rm -rf lk ~/.cache/lk
