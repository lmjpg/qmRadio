build: clean ui.go
	go build -ldflags "-s -w"
	upx ./qmRadio

run: ui.go
	go run .

ui.go:
	miqt-uic -Qt6 -InFile ui/qmRadio.ui -OutFile ui.go
	sed -i -e 's/qt.Orientation__Horizontal/qt.Horizontal/g' ui.go

clean:
	rm ./ui.go
	rm ./qmRadio