package tab

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func GetTabCards(tblCard *widget.Table, buttonSynchronization *widget.Button, cardAdd *widget.Button, cardDelete *widget.Button,
	cardUpdate *widget.Button, labelAlertCard *widget.Label) *container.TabItem {
	bottomContainer := container.New(layout.NewHBoxLayout(), cardAdd, cardDelete, cardUpdate, labelAlertCard)
	containerTblCard := layout.NewBorderLayout(buttonSynchronization, bottomContainer, nil, nil)
	boxCard := container.New(containerTblCard, buttonSynchronization, tblCard, bottomContainer)
	return container.NewTabItem("Банковские карты", boxCard)
}
