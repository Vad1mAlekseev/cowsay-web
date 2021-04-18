.PHONY: build
build:
	go build -o build/ ./...
	cp web build -r

.PHONY: clean
clean:
	rm -rf build
