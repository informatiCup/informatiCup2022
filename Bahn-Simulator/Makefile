all: release

clean:
	rm -rf build-win build-linux build-darwin dist

dist: clean
	mkdir build-win
	mkdir build-linux
	mkdir build-darwin
	CGO_ENABLED=0 GOOS=windows go build -o build-win/bahn-simulator.exe
	CGO_ENABLED=0 GOOS=linux go build -o build-linux/bahn-simulator
	CGO_ENABLED=0 GOOS=darwin go build -o build-darwin/bahn-simulator

release: dist
	mkdir dist
	cp -r licenses-third-party build-win
	cp -r licenses-third-party build-linux
	cp -r licenses-third-party build-darwin
	cp LICENSE build-win
	cp LICENSE build-linux
	cp LICENSE build-darwin
	zip dist/bahn-simulator-win.zip build-win/* build-win/*/*
	tar czvf dist/bahn-simulator-linux.tar.gz build-linux/*
	tar czvf dist/bahn-simulator-mac.tar.gz build-darwin/*
