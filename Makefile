.PHONY: test release install build clean

test:
	$(MAKE) -C dogbag $@

release: test
	(cd dogbag && godocdown -signature . > README.markdown) || false
	cp dogbag/README.markdown .

install: test
	$(MAKE) -C dogbag $@
	go install

build:
	$(MAKE) -C dogbag $@

clean:
	$(MAKE) -C dogbag $@
