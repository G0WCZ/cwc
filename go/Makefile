build: */*.go */*/*.go
	go build ./...

test:
	go test ./...

dist: dist-pi dist-linux dist-osx dist-windows

dist-pi: dist/pi/reflector/cwc-reflector dist/pi/station/cwc-station

dist/pi/reflector/cwc-reflector:
	mkdir -p dist/pi/reflector
	GOOS=linux GOARCH=arm GOARM=5 go build -o dist/pi/reflector/cwc-reflector ./cwc-reflector

dist/pi/station/cwc-station:
	mkdir -p dist/pi/station
	GOOS=linux GOARCH=arm GOARM=5 go build -o dist/pi/station/cwc-station ./cwc-station

dist-linux: dist/linux/reflector/cwc-reflector dist/linux/station/cwc-station

dist/linux/reflector/cwc-reflector:
	mkdir -p dist/linux/reflector
	GOOS=linux GOARCH=amd64 go build -o dist/linux/reflector/cwc-reflector ./cwc-reflector

dist/linux/station/cwc-station:
	mkdir -p dist/linux/station
	GOOS=linux GOARCH=amd64 go build -o dist/linux/station/cwc-station ./cwc-station

dist-osx: dist/osx/reflector/cwc-reflector dist/osx/station/cwc-station

dist/osx/reflector/cwc-reflector:
	mkdir -p dist/osx/reflector
	GOOS=darwin GOARCH=amd64 go build -o dist/osx/reflector/cwc-reflector ./cwc-reflector

dist/osx/station/cwc-station:
	mkdir -p dist/osx/station
	GOOS=darwin GOARCH=amd64 go build -o dist/osx/station/cwc-reflector ./cwc-station

dist-windows: dist/windows/reflector/cwc-reflector dist/windows/cwc-station

dist/windows/reflector/cwc-reflector:
	mkdir -p dist/windows/reflector
	GOOS=windows GOARCH=386 go build -o dist/windows/reflector/cwc-reflector.exe ./cwc-reflector

dist/windows/reflector/cwc-station:
	mkdir -p dist/windows/station
	GOOS=windows GOARCH=386 go build -o dist/windows/station/cwc-station.exe ./cwc-station


