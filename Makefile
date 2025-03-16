build:
	go build -o bin/go-lang

run:
	./bin/go-lang

test:
	go test -v ./...

clean:
	rm -rf bin/*

.PHONY: build run test clean
