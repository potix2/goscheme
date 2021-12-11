.PHONY: build
build:
	go build goscheme.go
all: build

.PHONY: clean
clean:
	rm -f goscheme

.PHONY: test
test:
	go test -v ./...


