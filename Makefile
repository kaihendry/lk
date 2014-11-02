lk: main.go
	go build

dist:
	gox

test: lk
	open http://0.0.0.0:3000
	./lk

clean:
	rm -f lk
