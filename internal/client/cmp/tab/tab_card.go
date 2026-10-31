package tab

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/xChygyNx/gophkeeper/internal/client/cmp/function"
)

func GetTabCards(myTabs *function.Tabs, myButtons *function.Buttons,
	myLabels *function.Labels) *container.TabItem {
	bottomContainer := container.New(layout.NewHBoxLayout(), myButtons.ButtonCard, myButtons.ButtonCardDelete,
		myButtons.ButtonCardUpdate, myLabels.LabelAlertCard)
	containerTblCard := layout.NewBorderLayout(myButtons.ButtonTopSynchronization, bottomContainer, nil, nil)
	boxCard := container.New(containerTblCard, myButtons.ButtonTopSynchronization, myTabs.TblCard, bottomContainer)
	return container.NewTabItem("Банковские карты", boxCard)
}
