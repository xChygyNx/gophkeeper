package tab

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func GetTabTexts(tblText *widget.Table, buttonSynchronization *widget.Button, textAdd *widget.Button, textDelete *widget.Button,
	textUpdate *widget.Button, labelAlertText *widget.Label) *container.TabItem {
	bottomContainer := container.New(layout.NewHBoxLayout(), textAdd, textDelete, textUpdate, labelAlertText)
	containerTblText := layout.NewBorderLayout(buttonSynchronization, bottomContainer, nil, nil)
	boxText := container.New(containerTblText, buttonSynchronization, tblText, bottomContainer)
	return container.NewTabItem("Текстовые данные", boxText)
}
