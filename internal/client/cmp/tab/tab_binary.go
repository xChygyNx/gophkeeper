package tab

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func GetTabBinaries(tblBinary *widget.Table, buttonSynchronization *widget.Button, binaryAdd *widget.Button, binaryDelete *widget.Button,
	binaryDownload *widget.Button, labelAlertBinary *widget.Label) *container.TabItem {
	bottomContainer := container.New(layout.NewHBoxLayout(), binaryAdd, binaryDelete, binaryDownload, labelAlertBinary)
	containerTblBinary := layout.NewBorderLayout(buttonSynchronization, bottomContainer, nil, nil)
	boxBinary := container.New(containerTblBinary, buttonSynchronization, tblBinary, bottomContainer)
	return container.NewTabItem("Файлы", boxBinary)
}
