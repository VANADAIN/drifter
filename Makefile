build:
	go build -o bin/drifter

run: build
	./bin/drifter

test:
	go test -v ./...