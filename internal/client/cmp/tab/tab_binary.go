package tab

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"github.com/xChygyNx/gophkeeper/internal/client/gui_elements"
)

func GetTabBinaries(myTabs *gui_elements.Tabs, myButtons *gui_elements.Buttons,
	myLabels *gui_elements.Labels) *container.TabItem {
	bottomContainer := container.New(layout.NewHBoxLayout(), myButtons.ButtonBinaryUpload, myButtons.ButtonBinaryDelete,
		myButtons.ButtonBinaryDownload, myLabels.LabelAlertBinary)
	containerTblBinary := layout.NewBorderLayout(myButtons.ButtonTopSynchronization, bottomContainer, nil, nil)
	boxBinary := container.New(containerTblBinary, myButtons.ButtonTopSynchronization, myTabs.TblBinary, bottomContainer)
	return container.NewTabItem("Файлы", boxBinary)
}
