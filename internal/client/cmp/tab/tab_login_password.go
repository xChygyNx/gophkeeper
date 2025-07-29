package tab

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func GetTabLoginPassword(tblLoginPassword *widget.Table, buttonSynchronization *widget.Button, loginPasswordAdd *widget.Button,
	loginPasswordDelete *widget.Button, loginPasswordUpdate *widget.Button, labelAlertLoginPassword *widget.Label) *container.TabItem {
	bottomContainer := container.New(layout.NewHBoxLayout(), loginPasswordAdd, loginPasswordDelete, loginPasswordUpdate, labelAlertLoginPassword)
	containerTblLoginPassword := layout.NewBorderLayout(buttonSynchronization, bottomContainer, nil, nil)
	boxLoginPassword := container.New(containerTblLoginPassword, buttonSynchronization, tblLoginPassword, bottomContainer)
	return container.NewTabItem("Логин пароль", boxLoginPassword)
}
