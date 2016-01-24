gitVersion := $(shell git rev-parse --short HEAD)

lk: main.go thumb.go
	time go build -ldflags "-X main.gitVersion $(gitVersion)"

docker: lk
	docker build -t lk .
	docker run -it -p 3000:3000 --rm lk

clean:
	rm -rf lk
