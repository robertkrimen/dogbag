.PHONY: test release build clean

test:
	go test -i
	go test -v
	$(MAKE) build

release: test
	godocdown -signature . > README.markdown

build:
	go build -v

clean:
	rm -f zipbag
	rm -rf test.*
