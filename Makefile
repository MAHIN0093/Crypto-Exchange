run:
	go build -o bin/go-lang && ./bin/go-lang

test:
	go test -v ./...

clean:
	rm -rf bin/*

.PHONY: build run test clean
