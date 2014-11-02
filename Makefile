lk: main.go
	go build
	xdg-open http://0.0.0.0:3000
	./lk

dist:
	gox

clean:
	rm -f lk
