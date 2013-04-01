.PHONY: test release build clean install

test: zipbag_template.go
	$(MAKE) -C zipbag $@
	go test -i
	go test -v
	$(MAKE) build
	./dogbag -output /dev/null ./test/assets 

release:
	$(MAKE) -C .. $@

install: test
	go install -v

build:
	go build -v

clean:
	rm -f dogbag assets_dogbag.go
	$(MAKE) -C zipbag $@

zipbag_template.go: zipbag_template/main.go zipbag/zipbag.go
	go run $< < $(word 2, $^) > $@
