build:
	go build

build-linux:
	GOOS=linux go build -o kubetop-linux

docker: build-linux
	docker build -t kubetop .

install-libs:
	govendor install +vendor,^program

clean:
	rm kubetop kubetop-linux 2>/dev/null || true
.PHONY: build build-linux docker install-libs clean
