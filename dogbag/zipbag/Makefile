.PHONY: test release build clean

test:
	go test -i
	go test
	$(MAKE) build

release:
	$(MAKE) -C .. $@

build:
	go build

clean:
	rm -f zipbag
	rm -rf test.*
