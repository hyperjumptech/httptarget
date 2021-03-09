GO111MODULE=on

.PHONY: clean build-all

build-linux:
	mkdir -p build/linux
	env GOOS=linux GOARCH=amd64 go build -o build/linux ./...
	mv build/linux/cmd build/linux/httptarget.app
	chmod 755 build/linux/httptarget.app
	cd build/linux; tar -zcvf httptarget-linux.tar.gz httptarget.app; rm -f httptarget.app

build-windows:
	mkdir -p build/windows
	env GOOS=windows GOARCH=amd64 go build -o build/windows ./...
	mv build/windows/cmd.exe build/windows/httptarget.exe
	cd build/windows; tar -zcvf httptarget-windows.tar.gz httptarget.exe; rm -f httptarget.exe

build-macos:
	mkdir -p build/macos
	env GOOS=darwin GOARCH=amd64 go build -o build/macos ./...
	mv build/macos/cmd build/macos/httptarget.app
	chmod 755 build/macos/httptarget.app
	cd build/macos; tar -zcvf httptarget-macos.tar.gz httptarget.app; rm -f httptarget.app

build-all: build-linux build-windows build-macos
	echo "Done"

clean:
	rm -Rf build