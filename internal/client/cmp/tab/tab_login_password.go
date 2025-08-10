package tab

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/xChygyNx/gophkeeper/internal/client/cmp/function"
)

func GetTabLoginPassword(myTabs *function.Tabs, myButtons *function.Buttons,
	myLabels *function.Labels) *container.TabItem {
	bottomContainer := container.New(layout.NewHBoxLayout(), myButtons.ButtonLoginPassword, myButtons.ButtonLoginPasswordDelete,
		myButtons.ButtonLoginPasswordUpdate, myLabels.LabelAlertLoginPassword)
	containerTblLoginPassword := layout.NewBorderLayout(myButtons.ButtonTopSynchronization, bottomContainer, nil, nil)
	boxLoginPassword := container.New(containerTblLoginPassword, myButtons.ButtonTopSynchronization, myTabs.TblLoginPassword, bottomContainer)
	return container.NewTabItem("Логин пароль", boxLoginPassword)
}
