.PHONY: build run clean

build: clean ui.go popup.go
	go build -ldflags "-s -w"
	upx ./qmRadio

run: ui.go popup.go
	go run .

ui.go:
	miqt-uic -Qt6 -InFile ui/qmRadio.ui -OutFile ui.go
	sed -i -e 's/qt.Orientation__Horizontal/qt.Horizontal/g' ui.go

popup.go:
	miqt-uic -Qt6 -InFile ui/popup.ui -OutFile popup.go
	sed -i -e 's/qt.AlignmentFlag__/qt./g' -e 's/qt.Orientation__Horizontal/qt.Horizontal/g' -e 's/QDialogButtonBox__StandardButton__/QDialogButtonBox__/g' popup.go

clean:
	rm ./ui.go
	rm ./popup.go
	rm ./qmRadio