package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/xChygyNx/gophkeeper/internal/client/api/events"
	"github.com/xChygyNx/gophkeeper/internal/client/cmp/form"
	"github.com/xChygyNx/gophkeeper/internal/client/cmp/function"
	"github.com/xChygyNx/gophkeeper/internal/client/cmp/tab"
	"github.com/xChygyNx/gophkeeper/internal/client/consts"
	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/table"
	"github.com/xChygyNx/gophkeeper/internal/client/storage/errors"
	"github.com/xChygyNx/gophkeeper/internal/client/storage/labels"
	"github.com/xChygyNx/gophkeeper/internal/client/storage/windows"
)

func InitGUI(log *logrus.Logger, application fyne.App, client *events.Event) {
	window := application.NewWindow("gophkeeper")

	window.Resize(fyne.NewSize(windows.WindowSwitcherWidth.Size(), windows.WindowSwitcherHeight.Size()))
	var dataTblLoginPassword = [][]string{{labels.NameItem, labels.DescriptionItem, labels.LoginItem, labels.PasswordItem,
		labels.CreatedAtItem, labels.UpdatedAtItem}}
	var dataTblText = [][]string{{labels.NameItem, labels.DescriptionItem, labels.DataItem, labels.CreatedAtItem, labels.UpdatedAtItem}}
	var dataTblCard = [][]string{{labels.NameItem, labels.DescriptionItem, labels.PaymentSystemItem, labels.NumberItem, labels.HolderItem, labels.CVCItem,
		labels.EndDateItem, labels.CreatedAtItem, labels.UpdatedAtItem}}
	var dataTblBinary = [][]string{{labels.NameItem, labels.CreatedAtItem}}

	var indexTblLoginPassword = 0
	var selectedRowTblLoginPassword []string
	var indexTblText = 0
	var selectedRowTblText []string
	var indexTblCard = 0
	var selectedRowTblCard []string
	var indexTblBinary = 0
	var selectedRowTblBinary []string

	var radioOptions = []string{labels.RadioBtnLogin, labels.RadioBtnRegistration}
	var accessToken = model.Token{}
	var password string
	var exist bool
	var err error
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
	var buttonAuth *widget.Button

	var buttonTopBack *widget.Button
	var buttonTopSynchronization *widget.Button

	var buttonLoginPassword *widget.Button
	var buttonText *widget.Button
	var buttonCard *widget.Button

	var buttonLoginPasswordCreate *widget.Button
	var buttonLoginPasswordDelete *widget.Button
	var buttonLoginPasswordUpdate *widget.Button
	var buttonLoginPasswordFormUpdate *widget.Button

	var buttonTextCreate *widget.Button
	var buttonTextDelete *widget.Button
	var buttonTextUpdate *widget.Button
	var buttonTextFormUpdate *widget.Button

	var buttonCardCreate *widget.Button
	var buttonCardDelete *widget.Button
	var buttonCardUpdate *widget.Button
	var buttonCardFormUpdate *widget.Button

	var buttonBinaryUpload *widget.Button
	var buttonBinaryDelete *widget.Button
	var buttonBinaryDownload *widget.Button
	//---------------------------------------------------------------------- tabs
	var containerTabs *container.AppTabs
	var tblLoginPassword *widget.Table
	var tblText *widget.Table
	var tblCard *widget.Table
	var tblBinary *widget.Table
	var tabLoginPassword *container.TabItem
	var tabText *container.TabItem
	var tabCard *container.TabItem
	var tabBinary *container.TabItem
	//---------------------------------------------------------------------- entries init
	separator := widget.NewSeparator()
	usernameLoginEntry := widget.NewEntry()
	passwordLoginEntry := widget.NewPasswordEntry()
	usernameRegistrationEntry := widget.NewEntry()
	passwordRegistrationEntry := widget.NewPasswordEntry()
	passwordConfirmationRegistrationEntry := widget.NewPasswordEntry()

	textNameEntryCreate := widget.NewEntry()
	textDescriptionEntryCreate := widget.NewEntry()
	textEntryCreate := widget.NewEntry()

	textNameEntryUpdate := widget.NewEntry()
	textNameEntryUpdate.Disable()
	textDescriptionEntryUpdate := widget.NewEntry()
	textDescriptionEntryUpdate.Disable()
	textEntryUpdate := widget.NewEntry()

	loginPasswordNameEntryCreate := widget.NewEntry()
	loginPasswordDescriptionEntryCreate := widget.NewEntry()
	loginEntryCreate := widget.NewEntry()
	passwordEntryCreate := widget.NewEntry()

	loginPasswordNameEntryUpdate := widget.NewEntry()
	loginPasswordNameEntryUpdate.Disable()
	loginPasswordDescriptionEntryUpdate := widget.NewEntry()
	loginPasswordDescriptionEntryUpdate.Disable()
	loginEntryUpdate := widget.NewEntry()
	passwordEntryUpdate := widget.NewEntry()

	cardNameEntryCreate := widget.NewEntry()
	cardDescriptionEntryCreate := widget.NewEntry()
	paymentSystemEntryCreate := widget.NewEntry()
	numberEntryCreate := widget.NewEntry()
	holderEntryCreate := widget.NewEntry()
	endDateEntryCreate := widget.NewEntry()
	cvcEntryCreate := widget.NewEntry()

	cardNameEntryUpdate := widget.NewEntry()
	cardNameEntryUpdate.Disable()
	cardDescriptionEntryUpdate := widget.NewEntry()
	cardDescriptionEntryUpdate.Disable()
	paymentSystemEntryUpdate := widget.NewEntry()
	numberEntryUpdate := widget.NewEntry()
	holderEntryUpdate := widget.NewEntry()
	endDateEntryUpdate := widget.NewEntry()
	cvcEntryUpdate := widget.NewEntry()

	//---------------------------------------------------------------------- labels init
	labelAlertAuth := widget.NewLabel("")
	labelAlertLoginPassword := widget.NewLabel("")
	labelAlertLoginPasswordCreate := widget.NewLabel("")
	labelAlertLoginPasswordUpdate := widget.NewLabel("")
	labelAlertText := widget.NewLabel("")
	labelAlertTextCreate := widget.NewLabel("")
	labelAlertTextUpdate := widget.NewLabel("")
	labelAlertCard := widget.NewLabel("")
	labelAlertCardCreate := widget.NewLabel("")
	labelAlertCardUpdate := widget.NewLabel("")
	labelAlertBinary := widget.NewLabel("")

	labelAlertAuth.Hide()
	labelAlertLoginPassword.Hide()
	labelAlertLoginPasswordCreate.Hide()
	labelAlertLoginPasswordUpdate.Hide()
	labelAlertText.Hide()
	labelAlertTextCreate.Hide()
	labelAlertTextUpdate.Hide()
	labelAlertCard.Hide()
	labelAlertCardCreate.Hide()
	labelAlertCardUpdate.Hide()
	labelAlertBinary.Hide()
	//---------------------------------------------------------------------- forms init
	formLogin := form.GetFormLogin(usernameLoginEntry, passwordLoginEntry)
	formRegistration := form.GetFormRegistration(usernameRegistrationEntry, passwordRegistrationEntry, passwordConfirmationRegistrationEntry)

	formLoginPasswordCreate := form.GetFormLoginPassword(loginPasswordNameEntryCreate, loginPasswordDescriptionEntryCreate, loginEntryCreate, passwordEntryCreate)
	formLoginPasswordUpdate := form.GetFormLoginPassword(loginPasswordNameEntryUpdate, loginPasswordDescriptionEntryUpdate, loginEntryUpdate, passwordEntryUpdate)

	formTextCreate := form.GetFormText(textNameEntryCreate, textDescriptionEntryCreate, textEntryCreate)
	formTextUpdate := form.GetFormText(textNameEntryUpdate, textDescriptionEntryUpdate, textEntryUpdate)

	formCardCreate := form.GetFormCard(cardNameEntryCreate, cardDescriptionEntryCreate, paymentSystemEntryCreate, numberEntryCreate, holderEntryCreate, endDateEntryCreate, cvcEntryCreate)
	formCardUpdate := form.GetFormCard(cardNameEntryUpdate, cardDescriptionEntryUpdate, paymentSystemEntryUpdate, numberEntryUpdate, holderEntryUpdate, endDateEntryUpdate, cvcEntryUpdate)
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
	buttonTopSynchronization = widget.NewButton(labels.BtnUpdateData, func() {
		dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, err = client.Synchronization(password, accessToken)
		if err != nil {
			labelAlertText.Show()
			labelAlertCard.Show()
			labelAlertLoginPassword.Show()
			labelAlertBinary.Show()
			labelAlertText.SetText(errors.ErrSynchronization)
			labelAlertCard.SetText(errors.ErrSynchronization)
			labelAlertLoginPassword.SetText(errors.ErrSynchronization)
			labelAlertLoginPassword.SetText(errors.ErrSynchronization)
		} else {
			tblLoginPassword.Resize(fyne.NewSize(float32(len(dataTblLoginPassword)), float32(len(dataTblLoginPassword[0]))))
			tblLoginPassword.Refresh()
			tblText.Resize(fyne.NewSize(float32(len(dataTblText)), float32(len(dataTblText[0]))))
			tblText.Refresh()
			tblCard.Resize(fyne.NewSize(float32(len(dataTblCard)), float32(len(dataTblCard[0]))))
			tblCard.Refresh()
			window.SetContent(containerTabs)
		}
	})
	buttonLoginPassword = widget.NewButton(labels.BtnAddLoginPassword, func() {
		window.SetContent(containerFormLoginPasswordCreate)
		window.Show()
	})
	buttonText = widget.NewButton(labels.BtnAddText, func() {
		window.SetContent(containerFormTextCreate)
		window.Show()
	})
	buttonCard = widget.NewButton(labels.BtnAddCard, func() {
		window.SetContent(containerFormCardCreate)
		window.Show()
	})
	//---------------------------------------------------------------------- login password event delete
	buttonLoginPasswordDelete = widget.NewButton(labels.BtnDeleteLoginPassword, func() {
		function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)
		if indexTblLoginPassword > 0 && client.LoginPasswordDelete(selectedRowTblLoginPassword, accessToken) != nil {
			// Удаляем строку с индексом indexTblLoginPassword
			dataTblLoginPassword = table.RemoveRow(dataTblLoginPassword, indexTblLoginPassword)
			indexTblLoginPassword = 0
		} else {
			log.Error(errors.ErrLoginPasswordTblIndexDelete)
			labelAlertLoginPassword.Show()
			labelAlertLoginPassword.SetText(errors.ErrLoginPasswordTblIndexDelete)
		}
	})
	//---------------------------------------------------------------------- text event delete
	buttonTextDelete = widget.NewButton(labels.BtnDeleteText, func() {
		function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)
		if indexTblText > 0 && client.TextDelete(selectedRowTblText, accessToken) != nil {
			// Удаляем строку с индексом indexTblText
			dataTblText = table.RemoveRow(dataTblText, indexTblText)
			indexTblText = 0
		} else {
			log.Error(errors.ErrTextTblIndexDelete)
			labelAlertText.Show()
			labelAlertText.SetText(errors.ErrTextTblIndexDelete)
		}
	})
	//---------------------------------------------------------------------- card event delete
	buttonCardDelete = widget.NewButton(labels.BtnDeleteCard, func() {
		function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)
		if indexTblCard > 0 && client.CardDelete(selectedRowTblCard, accessToken) != nil {
			// Удаляем строку с индексом indexTblCard
			dataTblCard = table.RemoveRow(dataTblCard, indexTblCard)
			indexTblCard = 0
		} else {
			log.Error(errors.ErrCardTblIndexDelete)
			labelAlertCard.Show()
			labelAlertCard.SetText(errors.ErrCardTblIndexDelete)
		}
	})
	//---------------------------------------------------------------------- binary event delete
	buttonBinaryDelete = widget.NewButton(labels.BtnDeleteBinary, func() {
		function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)
		if indexTblBinary > 0 && client.FileRemove(selectedRowTblBinary, accessToken) != nil {
			// Удаляем строку с индексом indexTblBinary
			dataTblBinary = table.RemoveRow(dataTblBinary, indexTblBinary)
			indexTblBinary = 0
		} else {
			log.Error(errors.ErrBinaryTblIndexDelete)
			labelAlertBinary.Show()
			labelAlertBinary.SetText(errors.ErrBinaryTblIndexDelete)
		}
	})
	//---------------------------------------------------------------------- switch form update
	buttonLoginPasswordUpdate = widget.NewButton(labels.BtnUpdateLoginPassword, func() {
		if indexTblLoginPassword > 0 {
			function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)
			function.SetLoginPasswordData(selectedRowTblLoginPassword, loginPasswordNameEntryUpdate, loginPasswordDescriptionEntryUpdate, loginEntryUpdate, passwordEntryUpdate)
			window.SetContent(containerFormLoginPasswordUpdate)
			window.Show()
		} else {
			log.Error(errors.ErrLoginPasswordTblIndexUpdate)
			labelAlertLoginPassword.Show()
			labelAlertLoginPassword.SetText(errors.ErrLoginPasswordTblIndexUpdate)
		}
	})
	buttonTextUpdate = widget.NewButton(labels.BtnUpdateText, func() {
		if indexTblText > 0 {
			function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)
			function.SetTextData(selectedRowTblText, textNameEntryUpdate, textDescriptionEntryUpdate, textEntryUpdate)
			window.SetContent(containerFormTextUpdate)
			window.Show()
		} else {
			log.Error(errors.ErrTextTblIndexUpdate)
			labelAlertText.Show()
			labelAlertText.SetText(errors.ErrTextTblIndexUpdate)
		}
	})
	buttonCardUpdate = widget.NewButton(labels.BtnUpdateCard, func() {
		if indexTblCard > 0 {
			function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)
			function.SetCardData(selectedRowTblCard, cardNameEntryUpdate, cardDescriptionEntryUpdate, paymentSystemEntryUpdate,
				numberEntryUpdate, holderEntryUpdate, cvcEntryUpdate, endDateEntryUpdate)
			window.SetContent(containerFormCardUpdate)
			window.Show()
		} else {
			log.Error(errors.ErrCardTblIndexUpdate)
			labelAlertCard.Show()
			labelAlertCard.SetText(errors.ErrCardTblIndexUpdate)
		}
	})
	//---------------------------------------------------------------------- login password event update
	buttonLoginPasswordFormUpdate = widget.NewButton(labels.BtnUpdate, func() {
		labelAlertLoginPasswordUpdate.Show()
		function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)

		errMsg, valid := function.ValidateLoginPasswordForm(loginPasswordNameEntryUpdate, loginPasswordDescriptionEntryUpdate,
			loginEntryUpdate, passwordEntryUpdate)
		if valid {
			err = client.LoginPasswordUpdate(loginPasswordNameEntryUpdate.Text, password,
				loginEntryUpdate.Text, passwordEntryUpdate.Text, accessToken)
			if err != nil {
				labelAlertLoginPasswordUpdate.SetText(errors.ErrLoginPasswordUpdate)
				log.Error(err)
			} else {

				dataTblLoginPassword = table.UpdateRowLoginPassword(loginEntryUpdate.Text, passwordEntryUpdate.Text, dataTblLoginPassword,
					indexTblLoginPassword)
				log.Info("Логин-пароль изменен")

				labelAlertLoginPasswordUpdate.Hide()
				formLoginPasswordUpdate.Refresh()
				window.SetContent(containerTabs)
				window.Show()
			}
		} else {
			labelAlertLoginPasswordUpdate.SetText(errMsg)
			log.Error(errMsg)
		}
	})
	//---------------------------------------------------------------------- text event update
	buttonTextFormUpdate = widget.NewButton(labels.BtnUpdate, func() {
		labelAlertTextUpdate.Show()
		function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)

		errMsg, valid := function.ValidateTextForm(textNameEntryUpdate, textDescriptionEntryUpdate, textEntryUpdate)
		if valid {
			err = client.TextUpdate(textNameEntryUpdate.Text, password, textEntryUpdate.Text, accessToken)
			if err != nil {
				labelAlertTextUpdate.SetText(errors.ErrTextUpdate)
				log.Error(err)
			} else {

				dataTblText = table.UpdateRowText(textEntryUpdate.Text, dataTblText, indexTblText)
				log.Info("Текст изменен")

				labelAlertTextUpdate.Hide()
				formLoginPasswordUpdate.Refresh()
				window.SetContent(containerTabs)
				window.Show()
			}
		} else {
			labelAlertTextUpdate.SetText(errMsg)
			log.Error(errMsg)
		}
	})
	//---------------------------------------------------------------------- card event update
	buttonCardFormUpdate = widget.NewButton(labels.BtnUpdate, func() {
		labelAlertCardUpdate.Show()
		function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)

		errMsg, valid := function.ValidateCardForm(cardNameEntryUpdate, cardDescriptionEntryUpdate, paymentSystemEntryUpdate,
			numberEntryUpdate, holderEntryUpdate, cvcEntryUpdate, endDateEntryUpdate)
		if valid {
			err = client.CardUpdate(cardNameEntryUpdate.Text, password, paymentSystemEntryUpdate.Text, numberEntryUpdate.Text,
				holderEntryUpdate.Text, cvcEntryUpdate.Text, endDateEntryUpdate.Text, accessToken)
			if err != nil {
				labelAlertCardUpdate.SetText(errors.ErrCardUpdate)
				log.Error(err)
			} else {

				dataTblCard = table.UpdateRowCard(paymentSystemEntryUpdate.Text, numberEntryUpdate.Text, holderEntryUpdate.Text,
					cvcEntryUpdate.Text, endDateEntryUpdate.Text, dataTblCard, indexTblCard)
				log.Info("Карта изменена")

				labelAlertCardUpdate.Hide()
				formCardUpdate.Refresh()
				window.SetContent(containerTabs)
				window.Show()
			}
		} else {
			labelAlertCardUpdate.SetText(errMsg)
			log.Error(errMsg)
		}
	})
	//----------------------------------------------------------------------  upload event
	buttonBinaryUpload = widget.NewButton(labels.BtnUploadBinary, func() {
		fileDialog := dialog.NewFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					labelAlertBinary.Show()
					exist = table.SearchByColumn(dataTblBinary, 0, r.URI().Name()) // search in map
					if exist {
						labelAlertBinary.SetText(errors.ErrBinaryExist)
						log.Error(labelAlertBinary.Text)
					}

					data, err := io.ReadAll(r)
					if len(data) > client.GetConfig().FileSize {
						labelAlertBinary.SetText(errors.ErrFileSize + "Размер загружаемого файла: " + strconv.Itoa(len(data)) + " байт")
						log.Error(errors.ErrFileSize + "Размер загружаемого файла: " + strconv.Itoa(len(data)) + " байт")
					} else if err != nil {
						labelAlertBinary.SetText(errors.ErrUpload)
						log.Error(err)
					} else {
						name, err := client.FileUpload(r.URI().Name(), password, data, accessToken)
						if err != nil {
							labelAlertBinary.SetText(errors.ErrUpload)
							log.Error(err)
						} else {
							dataTblBinary = append(dataTblBinary, []string{name, time.Now().Format(consts.DateAndTimeFormat)})
							log.Info("Файл добавлен: " + name)

							labelAlertBinary.Hide()
							window.SetContent(containerTabs)
							window.Show()
						}
					}
				}
			}, window)
		fileDialog.Show()
	})
	//----------------------------------------------------------------------  download event
	buttonBinaryDownload = widget.NewButton(labels.BtnDownloadBinary, func() {
		if indexTblBinary > 0 {
			err = client.FileDownload(selectedRowTblBinary[0], password, accessToken)
			if err != nil {
				labelAlertBinary.SetText(errors.ErrLogin)
				log.Error(err)
			}
			log.Info(indexTblBinary)
			log.Info(selectedRowTblBinary[0])
		} else {
			log.Error(errors.ErrBinaryTblIndexDownload)
			labelAlertBinary.Show()
			labelAlertBinary.SetText(errors.ErrBinaryTblIndexDownload)
		}
	})
	//----------------------------------------------------------------------
	buttonTopBack = widget.NewButton(labels.BtnBack, func() {
		function.ClearLoginPassword(loginPasswordNameEntryCreate, loginPasswordDescriptionEntryCreate, loginEntryCreate, passwordEntryCreate)
		function.ClearText(textNameEntryCreate, textDescriptionEntryCreate, textEntryCreate)
		function.ClearCard(cardNameEntryCreate, cardDescriptionEntryCreate, paymentSystemEntryCreate, numberEntryCreate, holderEntryCreate, endDateEntryCreate, cvcEntryCreate)
		function.ClearLoginPassword(loginPasswordNameEntryUpdate, loginPasswordDescriptionEntryUpdate, loginEntryUpdate, passwordEntryUpdate)
		function.ClearText(textNameEntryUpdate, textDescriptionEntryUpdate, textEntryUpdate)
		function.ClearCard(cardNameEntryUpdate, cardDescriptionEntryUpdate, paymentSystemEntryUpdate, numberEntryUpdate, holderEntryUpdate, endDateEntryUpdate, cvcEntryUpdate)
		labelAlertLoginPasswordCreate.Hide()
		labelAlertLoginPasswordUpdate.Hide()
		labelAlertTextCreate.Hide()
		labelAlertTextUpdate.Hide()
		labelAlertCardCreate.Hide()
		labelAlertCardUpdate.Hide()

		window.SetContent(containerTabs)
		window.Resize(fyne.NewSize(windows.WindowMainWidth.Size(), windows.WindowMainHeight.Size()))
		window.Show()
	})
	//---------------------------------------------------------------------- table login password init
	tblLoginPassword = widget.NewTable(
		func() (int, int) {
			return len(dataTblLoginPassword), len(dataTblLoginPassword[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(labels.TblLabel)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(dataTblLoginPassword[i.Row][i.Col])
		})
	function.SetDefaultColumnsWidthLoginPassword(tblLoginPassword)
	//---------------------------------------------------------------------- table text init
	tblText = widget.NewTable(
		func() (int, int) {
			return len(dataTblText), len(dataTblText[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(labels.TblLabel)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(dataTblText[i.Row][i.Col])
		})
	function.SetDefaultColumnsWidthText(tblText)
	//---------------------------------------------------------------------- table card init
	tblCard = widget.NewTable(
		func() (int, int) {
			return len(dataTblCard), len(dataTblCard[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(labels.TblLabel)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(dataTblCard[i.Row][i.Col])
		})
	function.SetDefaultColumnsWidthCard(tblCard)
	//---------------------------------------------------------------------- table binary init
	tblBinary = widget.NewTable(
		func() (int, int) {
			return len(dataTblBinary), len(dataTblBinary[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(labels.TblLabel)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(dataTblBinary[i.Row][i.Col])
		})
	function.SetDefaultColumnsWidthBinary(tblBinary)
	//---------------------------------------------------------------------- containerTabs
	tabLoginPassword = tab.GetTabLoginPassword(tblLoginPassword, buttonTopSynchronization, buttonLoginPassword, buttonLoginPasswordDelete, buttonLoginPasswordUpdate, labelAlertLoginPassword)
	tabText = tab.GetTabTexts(tblText, buttonTopSynchronization, buttonText, buttonTextDelete, buttonTextUpdate, labelAlertText)
	tabCard = tab.GetTabCards(tblCard, buttonTopSynchronization, buttonCard, buttonCardDelete, buttonCardUpdate, labelAlertCard)
	tabBinary = tab.GetTabBinaries(tblBinary, buttonTopSynchronization, buttonBinaryUpload, buttonBinaryDelete, buttonBinaryDownload, labelAlertBinary)
	containerTabs = container.NewAppTabs(tabLoginPassword, tabText, tabCard, tabBinary)
	//----------------------------------------------------------------------
	// Get selected row data
	tblLoginPassword.OnSelected = func(id widget.TableCellID) {
		indexTblLoginPassword = id.Row
		selectedRowTblLoginPassword = dataTblLoginPassword[id.Row]
	}
	tblText.OnSelected = func(id widget.TableCellID) {
		indexTblText = id.Row
		selectedRowTblText = dataTblText[id.Row]
	}
	tblCard.OnSelected = func(id widget.TableCellID) {
		indexTblCard = id.Row
		selectedRowTblCard = dataTblCard[id.Row]
	}
	tblBinary.OnSelected = func(id widget.TableCellID) {
		indexTblBinary = id.Row
		selectedRowTblBinary = dataTblBinary[id.Row]
	}
	//---------------------------------------------------------------------- auth event
	buttonAuth = widget.NewButton(labels.BtnSubmit, func() {
		labelAlertAuth.Show()
		if radioAuth.Selected == labels.RadioBtnLogin {
			errMsg, valid := function.ValidateLoginForm(usernameLoginEntry, passwordLoginEntry)
			if valid {
				accessToken, err = client.Authentication(usernameLoginEntry.Text, passwordLoginEntry.Text)
				if err != nil {
					labelAlertAuth.SetText(errors.ErrLogin)
					log.Error(err)
				} else {
					password = passwordLoginEntry.Text
					dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, err = client.Synchronization(password, accessToken)
					if err != nil {
						labelAlertAuth.SetText(errors.ErrLogin)
						log.Error(err)
					} else {
						window.SetContent(containerTabs)
						window.Resize(fyne.NewSize(windows.WindowMainWidth.Size(), windows.WindowMainHeight.Size()))
						window.Show()
					}
				}
			} else {
				labelAlertAuth.SetText(errMsg)
				log.Error(errMsg)
			}
		}
		if radioAuth.Selected == labels.RadioBtnRegistration {
			errMsg, valid := function.ValidateRegistrationForm(usernameRegistrationEntry, passwordRegistrationEntry, passwordConfirmationRegistrationEntry)
			if valid {
				exist, err = client.UserExist(usernameRegistrationEntry.Text)
				if err != nil {
					labelAlertAuth.SetText(errors.ErrRegistration)
					log.Error(err)
				}
				if exist {
					labelAlertAuth.SetText(errors.ErrUserExist)
					log.Error(errors.ErrUserExist)
				} else {
					accessToken, err = client.Registration(usernameRegistrationEntry.Text, passwordRegistrationEntry.Text)
					if err != nil {
						labelAlertAuth.SetText(errors.ErrRegistration)
						log.Error(err)
					} else {
						password = passwordRegistrationEntry.Text
						window.SetContent(containerTabs)
						window.Resize(fyne.NewSize(windows.WindowMainWidth.Size(), windows.WindowMainHeight.Size()))
						window.Show()
					}
				}
			} else {
				labelAlertAuth.SetText(errMsg)
				log.Error(errMsg)
			}
		}
	})
	//---------------------------------------------------------------------- login password event create
	buttonLoginPasswordCreate = widget.NewButton(labels.BtnAdd, func() {
		labelAlertLoginPasswordCreate.Show()
		function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)
		exist = table.SearchByColumn(dataTblLoginPassword, 0, loginPasswordNameEntryCreate.Text) // search in map
		if exist {
			labelAlertLoginPassword.SetText(errors.ErrLoginPasswordExist)
			log.Error(labelAlertLoginPassword.Text)
		}
		errMsg, valid := function.ValidateLoginPasswordForm(loginPasswordNameEntryCreate, loginPasswordDescriptionEntryCreate, loginEntryCreate,
			passwordEntryCreate)
		if valid {
			err = client.LoginPasswordCreate(loginPasswordNameEntryCreate.Text, loginPasswordDescriptionEntryCreate.Text, password,
				loginEntryCreate.Text, passwordEntryCreate.Text, accessToken)
			if err != nil {
				labelAlertLoginPasswordCreate.SetText(errors.ErrLoginPasswordCreate)
				log.Error(err)
			} else {
				dataTblLoginPassword = append(dataTblLoginPassword, []string{loginPasswordNameEntryCreate.Text, loginPasswordDescriptionEntryCreate.Text,
					loginEntryCreate.Text, passwordEntryCreate.Text, time.Now().Format(consts.DateAndTimeFormat), time.Now().Format(consts.DateAndTimeFormat)})

				function.ClearLoginPassword(loginPasswordNameEntryCreate, loginPasswordDescriptionEntryCreate, loginEntryCreate, passwordEntryCreate)
				log.Info("Логин-пароль добавлен")

				labelAlertLoginPasswordCreate.Hide()
				formLoginPasswordCreate.Refresh()
				window.SetContent(containerTabs)
				window.Show()
			}
		} else {
			labelAlertLoginPasswordCreate.SetText(errMsg)
			log.Error(errMsg)
		}
	})
	//---------------------------------------------------------------------- text event create
	buttonTextCreate = widget.NewButton(labels.BtnAdd, func() {
		labelAlertTextCreate.Show()
		function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)
		exist = table.SearchByColumn(dataTblText, 0, textNameEntryCreate.Text) // search in map
		if exist {
			labelAlertText.SetText(errors.ErrTextExist)
			log.Error(labelAlertText)
		}
		errMsg, valid := function.ValidateTextForm(textNameEntryCreate, textDescriptionEntryCreate, textEntryCreate)
		if valid {
			err = client.TextCreate(textNameEntryCreate.Text, textDescriptionEntryCreate.Text, password, textEntryCreate.Text, accessToken)
			if err != nil {
				labelAlertTextCreate.SetText(errors.ErrTextCreate)
				log.Error(err)
			} else {
				dataTblText = append(dataTblText, []string{textNameEntryCreate.Text, textDescriptionEntryCreate.Text, textEntryCreate.Text,
					time.Now().Format(consts.DateAndTimeFormat), time.Now().Format(consts.DateAndTimeFormat)})

				function.ClearText(textNameEntryCreate, textDescriptionEntryCreate, textEntryCreate)
				log.Info("Текст добавлен")

				labelAlertTextCreate.Hide()
				formTextCreate.Refresh()
				window.SetContent(containerTabs)
				window.Show()
			}
		} else {
			labelAlertTextCreate.SetText(errMsg)
			log.Error(errMsg)
		}
	})
	//---------------------------------------------------------------------- card event create
	buttonCardCreate = widget.NewButton(labels.BtnAdd, func() {
		labelAlertCardCreate.Show()
		function.HideLabelsTab(labelAlertLoginPassword, labelAlertText, labelAlertCard, labelAlertBinary)
		exist = table.SearchByColumn(dataTblCard, 0, cardNameEntryCreate.Text) // search in map
		if exist {
			labelAlertCard.SetText(errors.ErrCardExist)
			log.Print(labelAlertCard)
		}
		errMsg, valid := function.ValidateCardForm(cardNameEntryCreate, cardDescriptionEntryCreate, paymentSystemEntryCreate,
			numberEntryCreate, holderEntryCreate, cvcEntryCreate, endDateEntryCreate)
		if valid {
			err = client.CardCreate(cardNameEntryCreate.Text, cardDescriptionEntryCreate.Text, password,
				paymentSystemEntryCreate.Text, numberEntryCreate.Text, holderEntryCreate.Text,
				cvcEntryCreate.Text, endDateEntryCreate.Text, accessToken)
			if err != nil {
				labelAlertCardCreate.SetText(errors.ErrCardCreate)
				log.Error(err)
			} else {
				dataTblCard = append(dataTblCard, []string{cardNameEntryCreate.Text, cardDescriptionEntryCreate.Text,
					paymentSystemEntryCreate.Text, numberEntryCreate.Text, holderEntryCreate.Text,
					cvcEntryCreate.Text, endDateEntryCreate.Text, time.Now().Format(consts.DateAndTimeFormat),
					time.Now().Format(consts.DateAndTimeFormat)})

				function.ClearCard(cardNameEntryCreate, cardDescriptionEntryCreate, paymentSystemEntryCreate, numberEntryCreate, holderEntryCreate, endDateEntryCreate, cvcEntryCreate)
				log.Info("Карта добавлена")

				labelAlertCardCreate.Hide()
				formCardCreate.Refresh()
				window.SetContent(containerTabs)
				window.Show()
			}
		} else {
			labelAlertCardCreate.SetText(errMsg)
			log.Error(errMsg)
		}
	})
	//---------------------------------------------------------------------- containers init
	containerRadio = container.NewVBox(radioAuth)

	containerFormLogin = container.NewVBox(formLogin, buttonAuth, labelAlertAuth, separator, radioAuth)
	containerFormRegistration = container.NewVBox(formRegistration, buttonAuth, labelAlertAuth, separator, radioAuth)

	containerFormLoginPasswordCreate = container.NewVBox(buttonTopBack, formLoginPasswordCreate, buttonLoginPasswordCreate, labelAlertLoginPasswordCreate)
	containerFormLoginPasswordUpdate = container.NewVBox(buttonTopBack, formLoginPasswordUpdate, buttonLoginPasswordFormUpdate, labelAlertLoginPasswordUpdate)

	containerFormTextCreate = container.NewVBox(buttonTopBack, formTextCreate, buttonTextCreate, labelAlertTextCreate)
	containerFormTextUpdate = container.NewVBox(buttonTopBack, formTextUpdate, buttonTextFormUpdate, labelAlertTextUpdate)

	containerFormCardCreate = container.NewVBox(buttonTopBack, formCardCreate, buttonCardCreate, labelAlertCardCreate)
	containerFormCardUpdate = container.NewVBox(buttonTopBack, formCardUpdate, buttonCardFormUpdate, labelAlertCardUpdate)

	//----------------------------------------------------------------------
	window.SetContent(containerRadio)
	window.ShowAndRun()
}
