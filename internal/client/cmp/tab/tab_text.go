package tab

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"github.com/xChygyNx/gophkeeper/internal/client/gui_elements"
)

func GetTabTexts(myTabs *gui_elements.Tabs, myButtons *gui_elements.Buttons,
	myLabels *gui_elements.Labels) *container.TabItem {
	bottomContainer := container.New(layout.NewHBoxLayout(), myButtons.ButtonText, myButtons.ButtonTextDelete,
		myButtons.ButtonTextUpdate, myLabels.LabelAlertText)
	containerTblText := layout.NewBorderLayout(myButtons.ButtonTopSynchronization, bottomContainer, nil, nil)
	boxText := container.New(containerTblText, myButtons.ButtonTopSynchronization, myTabs.TblText, bottomContainer)
	return container.NewTabItem("Текстовые данные", boxText)
}
