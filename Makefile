lk: main.go
	go build
	xdg-open http://0.0.0.0:3000
	./lk

clean:
	rm -f lk
