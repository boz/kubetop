build:
	go build

install-libs:
	govendor install +vendor,^program

clean:
	rm kubetop 2>/dev/null || true
.PHONY: build install-libs clean
