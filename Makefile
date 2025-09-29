.PHONY: build run clean

build: clean ui.go popup.go
	go build -ldflags "-s -w"
	upx ./qmRadio

run: ui.go popup.go
	go run .

ui.go:
	miqt-uic -Qt6 -InFile ui/qmRadio.ui -OutFile ui.go
	sed -i \
		-e 's/qt.Orientation__Horizontal/qt.Horizontal/g' \
		-e 's/qt.PenStyle__NoPen/qt.NoPen/g' \
		-e 's/qt.QAbstractItemView__EditTrigger__/qt.QAbstractItemView__/g' \
		-e 's/qt.QAbstractItemView__SelectionBehavior__/qt.QAbstractItemView__/g' \
		-e 's/qt.QAbstractItemView__SelectionMode__/qt.QAbstractItemView__/g' \
		ui.go

popup.go:
	miqt-uic -Qt6 -InFile ui/popup.ui -OutFile popup.go
	sed -i
		-e 's/qt.AlignmentFlag__/qt./g' \
		-e 's/qt.Orientation__Horizontal/qt.Horizontal/g' \
		-e 's/QDialogButtonBox__StandardButton__/QDialogButtonBox__/g' popup.go

clean:
	rm -f ./ui.go
	rm -f ./popup.go
	rm -f ./qmRadio