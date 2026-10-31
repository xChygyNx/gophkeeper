package tab

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/xChygyNx/gophkeeper/internal/client/cmp/function"
)

func GetTabBinaries(myTabs *function.Tabs, myButtons *function.Buttons,
	myLabels *function.Labels) *container.TabItem {
	bottomContainer := container.New(layout.NewHBoxLayout(), myButtons.ButtonBinaryUpload, myButtons.ButtonBinaryDelete,
		myButtons.ButtonBinaryDownload, myLabels.LabelAlertBinary)
	containerTblBinary := layout.NewBorderLayout(myButtons.ButtonTopSynchronization, bottomContainer, nil, nil)
	boxBinary := container.New(containerTblBinary, myButtons.ButtonTopSynchronization, myTabs.TblBinary, bottomContainer)
	return container.NewTabItem("Файлы", boxBinary)
}
