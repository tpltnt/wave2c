main: wave2c.go
	go build wave2c.go
test:
	go test

valid: main
	./wave2c testdata/validinput.wav

.PHONY: clean
clean:
	rm *~
