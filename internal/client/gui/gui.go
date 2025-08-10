package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"

	"github.com/xChygyNx/gophkeeper/internal/client/api/events"
	"github.com/xChygyNx/gophkeeper/internal/client/cmp/function"
	"github.com/xChygyNx/gophkeeper/internal/client/cmp/tab"
	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/storage/labels"
	"github.com/xChygyNx/gophkeeper/internal/client/storage/windows"
)

func InitGUI(log *logrus.Logger, application fyne.App, client *events.Event) {
	window := application.NewWindow("gophkeeper")

	window.Resize(fyne.NewSize(windows.WindowSwitcherWidth.Size(), windows.WindowSwitcherHeight.Size()))
	dataTables := function.InitDataTables()
	myIndexes := function.InitIndexes()

	var radioOptions = []string{labels.RadioBtnLogin, labels.RadioBtnRegistration}
	var accessToken = model.Token{}
	var password string
	//---------------------------------------------------------------------- containers
	var containerRadio *fyne.Container

	var containerFormLogin *fyne.Container
	var containerFormRegistration *fyne.Container

	var containerFormLoginPasswordCreate *fyne.Container
	var containerFormTextCreate *fyne.Container
	var containerFormCardCreate *fyne.Container

	var containerFormLoginPasswordUpdate *fyne.Container
	var containerFormTextUpdate *fyne.Container
	var containerFormCardUpdate *fyne.Container
	//---------------------------------------------------------------------- buttons
	myButtons := function.GetButtons(client, window, log)
	//---------------------------------------------------------------------- tabs
	myTabs := function.GetTabs()
	//---------------------------------------------------------------------- entries init
	separator := widget.NewSeparator()
	myEntries := function.InitEntries()
	//---------------------------------------------------------------------- form validator init
	formValidator := function.NewFormValidator(myEntries)
	//---------------------------------------------------------------------- labels init
	myLabels := function.InitLabels()
	//---------------------------------------------------------------------- forms init
	myForms := function.InitForms(myEntries)
	//---------------------------------------------------------------------- radio event
	radioAuth := widget.NewRadioGroup(radioOptions, func(value string) {
		log.Println("Radio set to ", value)
		if value == labels.RadioBtnLogin {
			window.SetContent(containerFormLogin)
			window.Resize(fyne.NewSize(windows.WindowAuthWidth.Size(), windows.WindowAuthHeight.Size()))
			window.Show()
		}
		if value == labels.RadioBtnRegistration {
			window.SetContent(containerFormRegistration)
			window.Resize(fyne.NewSize(windows.WindowAuthWidth.Size(), windows.WindowAuthHeight.Size()))
			window.Show()
		}
	})
	//---------------------------------------------------------------------- buttons event
	myButtons.InitTopSynchronizationButton(myLabels, myTabs, dataTables, myTabs.ContainerTabs, password, accessToken)
	myButtons.InitLoginPasswordButton(containerFormLoginPasswordCreate)
	myButtons.InitAddTextButton(containerFormTextCreate)
	myButtons.InitAddCardButton(containerFormCardCreate)
	//---------------------------------------------------------------------- login password event delete
	myButtons.InitLoginPasswordDeleteButton(myLabels, dataTables, myIndexes, accessToken)
	//---------------------------------------------------------------------- text event delete
	myButtons.InitTextDeleteButton(myLabels, dataTables, myIndexes, accessToken)
	//---------------------------------------------------------------------- card event delete
	myButtons.InitCardDeleteButton(myLabels, dataTables, myIndexes, accessToken)
	//---------------------------------------------------------------------- binary event delete
	myButtons.InitBinaryDeleteButton(myLabels, dataTables, myIndexes, accessToken)
	//---------------------------------------------------------------------- switch form update
	myButtons.InitLoginPasswordUpdateButton(myLabels, myEntries, myIndexes, containerFormLoginPasswordUpdate)
	myButtons.InitTextUpdateButton(myLabels, myEntries, myIndexes, containerFormTextUpdate)
	myButtons.InitCardUpdateButton(myLabels, myEntries, myIndexes, containerFormCardUpdate)
	//---------------------------------------------------------------------- login password event update
	myButtons.InitLoginPasswordFormUpdateButton(myLabels, myEntries, myIndexes, dataTables, myForms, myTabs,
		formValidator, accessToken, password)
	//---------------------------------------------------------------------- text event update
	myButtons.InitTextFormUpdateButton(myLabels, myEntries, myIndexes, dataTables, myForms, myTabs,
		formValidator, accessToken, password)
	//---------------------------------------------------------------------- card event update
	myButtons.InitCardFormUpdateButton(myLabels, myEntries, myIndexes, dataTables, myForms, myTabs,
		formValidator, accessToken, password)
	//----------------------------------------------------------------------  upload event
	myButtons.InitBinaryUploadButton(myLabels, dataTables, myTabs, accessToken, password)
	//----------------------------------------------------------------------  download event
	myButtons.InitBinaryDownloadButton(myLabels, myIndexes, accessToken, password)
	//----------------------------------------------------------------------
	myButtons.InitTopBackButton(myLabels, myEntries, myTabs)
	//---------------------------------------------------------------------- table login password init
	myTabs.InitLoginPasswordTable(dataTables)
	//---------------------------------------------------------------------- table text init
	myTabs.InitTextTable(dataTables)
	//---------------------------------------------------------------------- table card init
	myTabs.InitCardTable(dataTables)
	//---------------------------------------------------------------------- table binary init
	myTabs.InitBinaryTable(dataTables)
	//---------------------------------------------------------------------- containerTabs
	tabLoginPassword := tab.GetTabLoginPassword(myTabs, myButtons, myLabels)
	tabText := tab.GetTabTexts(myTabs, myButtons, myLabels)
	tabCard := tab.GetTabCards(myTabs, myButtons, myLabels)
	tabBinary := tab.GetTabBinaries(myTabs, myButtons, myLabels)
	myTabs.ContainerTabs = container.NewAppTabs(tabLoginPassword, tabText, tabCard, tabBinary)
	//----------------------------------------------------------------------
	// Get selected row data
	myTabs.TblLoginPassword.OnSelected = func(id widget.TableCellID) {
		myIndexes.IndexTblLoginPassword = id.Row
		myIndexes.SelectedRowTblLoginPassword = dataTables.DataTblLoginPassword[id.Row]
	}
	myTabs.TblText.OnSelected = func(id widget.TableCellID) {
		myIndexes.IndexTblText = id.Row
		myIndexes.SelectedRowTblText = dataTables.DataTblText[id.Row]
	}
	myTabs.TblCard.OnSelected = func(id widget.TableCellID) {
		myIndexes.IndexTblCard = id.Row
		myIndexes.SelectedRowTblCard = dataTables.DataTblCard[id.Row]
	}
	myTabs.TblBinary.OnSelected = func(id widget.TableCellID) {
		myIndexes.IndexTblBinary = id.Row
		myIndexes.SelectedRowTblBinary = dataTables.DataTblBinary[id.Row]
	}
	//---------------------------------------------------------------------- auth event
	myButtons.InitAuthButton(myLabels, myEntries, myTabs, dataTables, radioAuth, formValidator, accessToken, password)

	//---------------------------------------------------------------------- login password event create
	myButtons.InitLoginPasswordCreateButton(myLabels, myEntries, myForms, myTabs, dataTables, formValidator,
		accessToken, password)
	//---------------------------------------------------------------------- text event create
	myButtons.InitTextCreateButton(myLabels, myEntries, myForms, myTabs, dataTables, formValidator,
		accessToken, password)
	//---------------------------------------------------------------------- card event create
	myButtons.InitCardCreateButton(myLabels, myEntries, myForms, myTabs, dataTables, formValidator,
		accessToken, password)
	//---------------------------------------------------------------------- containers init
	containerRadio = container.NewVBox(radioAuth)

	containerFormLogin = container.NewVBox(myForms.FormLogin, myButtons.ButtonAuth, myLabels.LabelAlertAuth, separator, radioAuth)
	containerFormRegistration = container.NewVBox(myForms.FormRegistration, myButtons.ButtonAuth, myLabels.LabelAlertAuth, separator, radioAuth)

	containerFormLoginPasswordCreate = container.NewVBox(myButtons.ButtonTopBack, myForms.FormLoginPasswordCreate, myButtons.ButtonLoginPasswordCreate, myLabels.LabelAlertLoginPasswordCreate)
	containerFormLoginPasswordUpdate = container.NewVBox(myButtons.ButtonTopBack, myForms.FormLoginPasswordUpdate, myButtons.ButtonLoginPasswordFormUpdate, myLabels.LabelAlertLoginPasswordUpdate)

	containerFormTextCreate = container.NewVBox(myButtons.ButtonTopBack, myForms.FormTextCreate, myButtons.ButtonTextCreate, myLabels.LabelAlertTextCreate)
	containerFormTextUpdate = container.NewVBox(myButtons.ButtonTopBack, myForms.FormTextUpdate, myButtons.ButtonTextFormUpdate, myLabels.LabelAlertTextUpdate)

	containerFormCardCreate = container.NewVBox(myButtons.ButtonTopBack, myForms.FormCardCreate, myButtons.ButtonCardCreate, myLabels.LabelAlertCardCreate)
	containerFormCardUpdate = container.NewVBox(myButtons.ButtonTopBack, myForms.FormCardUpdate, myButtons.ButtonCardFormUpdate, myLabels.LabelAlertCardUpdate)

	//----------------------------------------------------------------------
	window.SetContent(containerRadio)
	window.ShowAndRun()
}
