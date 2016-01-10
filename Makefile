gitVersion := $(shell git rev-parse --short HEAD)

lk: main.go thumb.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.gitVersion $(gitVersion)"

docker: lk
	docker build -t lk .
	docker run -it -p 3000:3000 --rm lk

clean:
	rm -rf lk
