.PHONY: build
build:
	go build -o build/ ./...
	cp web build -r
	mkdir build/configs && cp configs/* build/configs/

.PHONY: clean
clean:
	rm -rf build

.PHONY: before_commit
before_commit:
	go fmt ./...
	go test ./...
