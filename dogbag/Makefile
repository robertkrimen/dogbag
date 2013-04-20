.PHONY: test release build clean install

test:
	go build -o stub ./zipbag/stub
	./stub > zipbag/test_dogbag.go
	$(MAKE) -C zipbag $@
	$(MAKE) zipbag_template.go
	go test -i
	go test
	$(MAKE) build
	./dogbag -output /dev/null ./test/assets 

release:
	$(MAKE) -C .. $@

install: test
	go install

build:
	go build

clean:
	rm -f dogbag assets_dogbag.go stub
	$(MAKE) -C zipbag $@

zipbag_template.go: zipbag_template/main.go zipbag/zipbag.go
	go run $< < $(word 2, $^) > $@
