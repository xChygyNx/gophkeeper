package function

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/sirupsen/logrus"
	"github.com/xChygyNx/gophkeeper/internal/client/consts"
	"github.com/xChygyNx/gophkeeper/internal/client/storage/windows"
	"io"
	"strconv"
	"time"

	"github.com/xChygyNx/gophkeeper/internal/client/api/events"
	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/table"
	"github.com/xChygyNx/gophkeeper/internal/client/storage/errors"
	"github.com/xChygyNx/gophkeeper/internal/client/storage/labels"
)

type Buttons struct {
	ButtonAuth *widget.Button

	ButtonTopBack            *widget.Button
	ButtonTopSynchronization *widget.Button

	ButtonLoginPassword *widget.Button
	ButtonText          *widget.Button
	ButtonCard          *widget.Button

	ButtonLoginPasswordCreate     *widget.Button
	ButtonLoginPasswordDelete     *widget.Button
	ButtonLoginPasswordUpdate     *widget.Button
	ButtonLoginPasswordFormUpdate *widget.Button

	ButtonTextCreate     *widget.Button
	ButtonTextDelete     *widget.Button
	ButtonTextUpdate     *widget.Button
	ButtonTextFormUpdate *widget.Button

	ButtonCardCreate     *widget.Button
	ButtonCardDelete     *widget.Button
	ButtonCardUpdate     *widget.Button
	ButtonCardFormUpdate *widget.Button

	ButtonBinaryUpload   *widget.Button
	ButtonBinaryDelete   *widget.Button
	ButtonBinaryDownload *widget.Button

	client *events.Event
	window fyne.Window
	log    *logrus.Logger
}

func GetButtons(client *events.Event, window fyne.Window, log *logrus.Logger) *Buttons {
	return &Buttons{
		client: client,
		window: window,
		log:    log,
	}
}

func (b *Buttons) InitTopSynchronizationButton(myLabels *Labels, myTabs *Tabs, dataTables *DataTables,
	containerTabs *container.AppTabs, password string, accessToken model.Token) {
	b.ButtonTopSynchronization = widget.NewButton(labels.BtnUpdateData, func() {
		dataTblText, dataTblCard, dataTblLoginPassword, dataTblBinary, err := b.client.Synchronization(password, accessToken)
		dataTables.DataTblText = dataTblText
		dataTables.DataTblCard = dataTblCard
		dataTables.DataTblLoginPassword = dataTblLoginPassword
		dataTables.DataTblBinary = dataTblBinary
		if err != nil {
			myLabels.LabelAlertText.Show()
			myLabels.LabelAlertCard.Show()
			myLabels.LabelAlertLoginPassword.Show()
			myLabels.LabelAlertBinary.Show()
			myLabels.LabelAlertText.SetText(errors.ErrSynchronization)
			myLabels.LabelAlertCard.SetText(errors.ErrSynchronization)
			myLabels.LabelAlertLoginPassword.SetText(errors.ErrSynchronization)
			myLabels.LabelAlertLoginPassword.SetText(errors.ErrSynchronization)

		} else {
			myTabs.TblLoginPassword.Resize(fyne.NewSize(float32(len(dataTblLoginPassword)), float32(len(dataTblLoginPassword[0]))))
			myTabs.TblLoginPassword.Refresh()
			myTabs.TblText.Resize(fyne.NewSize(float32(len(dataTblText)), float32(len(dataTblText[0]))))
			myTabs.TblText.Refresh()
			myTabs.TblCard.Resize(fyne.NewSize(float32(len(dataTblCard)), float32(len(dataTblCard[0]))))
			myTabs.TblCard.Refresh()
			b.window.SetContent(containerTabs)
		}
	})
}

func (b *Buttons) InitLoginPasswordButton(containerFormLoginPasswordCreate *fyne.Container) {
	b.ButtonLoginPassword = widget.NewButton(labels.BtnAddLoginPassword, func() {
		b.window.SetContent(containerFormLoginPasswordCreate)
		b.window.Show()
	})
}

func (b *Buttons) InitAddTextButton(containerFormTextCreate *fyne.Container) {
	b.ButtonText = widget.NewButton(labels.BtnAddText, func() {
		b.window.SetContent(containerFormTextCreate)
		b.window.Show()
	})
}

func (b *Buttons) InitAddCardButton(containerFormCardCreate *fyne.Container) {
	b.ButtonCard = widget.NewButton(labels.BtnAddCard, func() {
		b.window.SetContent(containerFormCardCreate)
		b.window.Show()
	})
}

func (b *Buttons) InitLoginPasswordDeleteButton(myLabel *Labels, dataTables *DataTables, myIndexes *Indexes,
	accessToken model.Token) {
	b.ButtonLoginPasswordDelete = widget.NewButton(labels.BtnDeleteLoginPassword, func() {
		HideLabelsTab(myLabel.LabelAlertLoginPassword, myLabel.LabelAlertText,
			myLabel.LabelAlertCard, myLabel.LabelAlertBinary)
		if myIndexes.IndexTblLoginPassword > 0 && b.client.LoginPasswordDelete(myIndexes.SelectedRowTblLoginPassword, accessToken) != nil {
			// Удаляем строку с индексом indexTblLoginPassword
			dataTables.DataTblLoginPassword = table.RemoveRow(dataTables.DataTblLoginPassword, myIndexes.IndexTblLoginPassword)
			myIndexes.IndexTblLoginPassword = 0
		} else {
			b.log.Error(errors.ErrLoginPasswordTblIndexDelete)
			myLabel.LabelAlertLoginPassword.Show()
			myLabel.LabelAlertLoginPassword.SetText(errors.ErrLoginPasswordTblIndexDelete)
		}
	})
}

func (b *Buttons) InitTextDeleteButton(myLabels *Labels, dataTables *DataTables, myIndexes *Indexes,
	accessToken model.Token) {
	b.ButtonTextDelete = widget.NewButton(labels.BtnDeleteText, func() {
		HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText,
			myLabels.LabelAlertCard, myLabels.LabelAlertBinary)
		if myIndexes.IndexTblText > 0 && b.client.TextDelete(myIndexes.SelectedRowTblText, accessToken) != nil {
			// Удаляем строку с индексом indexTblText
			dataTables.DataTblText = table.RemoveRow(dataTables.DataTblText, myIndexes.IndexTblText)
			myIndexes.IndexTblText = 0
		} else {
			b.log.Error(errors.ErrTextTblIndexDelete)
			myLabels.LabelAlertText.Show()
			myLabels.LabelAlertText.SetText(errors.ErrTextTblIndexDelete)
		}
	})
}

func (b *Buttons) InitCardDeleteButton(myLabels *Labels, dataTables *DataTables, myIndexes *Indexes,
	accessToken model.Token) {
	b.ButtonCardDelete = widget.NewButton(labels.BtnDeleteCard, func() {
		HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText,
			myLabels.LabelAlertCard, myLabels.LabelAlertBinary)
		if myIndexes.IndexTblCard > 0 && b.client.CardDelete(myIndexes.SelectedRowTblCard, accessToken) != nil {
			// Удаляем строку с индексом indexTblCard
			dataTables.DataTblCard = table.RemoveRow(dataTables.DataTblCard, myIndexes.IndexTblCard)
			myIndexes.IndexTblCard = 0
		} else {
			b.log.Error(errors.ErrCardTblIndexDelete)
			myLabels.LabelAlertCard.Show()
			myLabels.LabelAlertCard.SetText(errors.ErrCardTblIndexDelete)
		}
	})
}

func (b *Buttons) InitBinaryDeleteButton(myLabels *Labels, dataTables *DataTables, myIndexes *Indexes,
	accessToken model.Token) {
	b.ButtonBinaryDelete = widget.NewButton(labels.BtnDeleteBinary, func() {
		HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText,
			myLabels.LabelAlertCard, myLabels.LabelAlertBinary)
		if myIndexes.IndexTblBinary > 0 && b.client.FileRemove(myIndexes.SelectedRowTblBinary, accessToken) != nil {
			// Удаляем строку с индексом indexTblBinary
			dataTables.DataTblBinary = table.RemoveRow(dataTables.DataTblBinary, myIndexes.IndexTblBinary)
			myIndexes.IndexTblBinary = 0
		} else {
			b.log.Error(errors.ErrBinaryTblIndexDelete)
			myLabels.LabelAlertBinary.Show()
			myLabels.LabelAlertBinary.SetText(errors.ErrBinaryTblIndexDelete)
		}
	})
}

func (b *Buttons) InitLoginPasswordUpdateButton(myLabels *Labels, myEntries *Entries, myIndexes *Indexes,
	containerForm *fyne.Container) {
	b.ButtonLoginPasswordUpdate = widget.NewButton(labels.BtnUpdateLoginPassword, func() {
		if myIndexes.IndexTblLoginPassword > 0 {
			HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText,
				myLabels.LabelAlertCard, myLabels.LabelAlertBinary)
			SetLoginPasswordData(myIndexes.SelectedRowTblLoginPassword, myEntries.LoginPasswordNameEntryUpdate,
				myEntries.LoginPasswordDescriptionEntryUpdate, myEntries.LoginEntryUpdate, myEntries.PasswordEntryUpdate)
			b.window.SetContent(containerForm)
			b.window.Show()
		} else {
			b.log.Error(errors.ErrLoginPasswordTblIndexUpdate)
			myLabels.LabelAlertLoginPassword.Show()
			myLabels.LabelAlertLoginPassword.SetText(errors.ErrLoginPasswordTblIndexUpdate)
		}
	})
}

func (b *Buttons) InitTextUpdateButton(myLabels *Labels, myEntries *Entries, myIndexes *Indexes,
	containerForm *fyne.Container) {
	b.ButtonTextUpdate = widget.NewButton(labels.BtnUpdateText, func() {
		if myIndexes.IndexTblText > 0 {
			HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText,
				myLabels.LabelAlertCard, myLabels.LabelAlertBinary)
			SetTextData(myIndexes.SelectedRowTblText, myEntries.TextNameEntryUpdate,
				myEntries.TextDescriptionEntryUpdate, myEntries.TextEntryUpdate)
			b.window.SetContent(containerForm)
			b.window.Show()
		} else {
			b.log.Error(errors.ErrTextTblIndexUpdate)
			myLabels.LabelAlertText.Show()
			myLabels.LabelAlertText.SetText(errors.ErrTextTblIndexUpdate)
		}
	})
}

func (b *Buttons) InitCardUpdateButton(myLabels *Labels, myEntries *Entries, myIndexes *Indexes,
	containerForm *fyne.Container) {
	b.ButtonCardUpdate = widget.NewButton(labels.BtnUpdateCard, func() {
		if myIndexes.IndexTblCard > 0 {
			HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText,
				myLabels.LabelAlertCard, myLabels.LabelAlertBinary)
			SetCardData(myIndexes.SelectedRowTblCard, myEntries.CardNameEntryUpdate,
				myEntries.CardDescriptionEntryUpdate, myEntries.PaymentSystemEntryUpdate, myEntries.NumberEntryUpdate,
				myEntries.HolderEntryUpdate, myEntries.CvcEntryUpdate, myEntries.EndDateEntryUpdate)
			b.window.SetContent(containerForm)
			b.window.Show()
		} else {
			b.log.Error(errors.ErrCardTblIndexUpdate)
			myLabels.LabelAlertCard.Show()
			myLabels.LabelAlertCard.SetText(errors.ErrCardTblIndexUpdate)
		}
	})
}

func (b *Buttons) InitLoginPasswordFormUpdateButton(myLabels *Labels, myEntries *Entries, myIndexes *Indexes,
	dataTables *DataTables, myForms *Forms, myTabs *Tabs, formValidator *FormValidator,
	accessToken model.Token, password string) {
	b.ButtonLoginPasswordFormUpdate = widget.NewButton(labels.BtnUpdate, func() {
		myLabels.LabelAlertLoginPasswordUpdate.Show()
		HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText, myLabels.LabelAlertCard,
			myLabels.LabelAlertBinary)

		errMsg, valid := formValidator.ValidateLoginPasswordForm(Update)
		if valid {
			err := b.client.LoginPasswordUpdate(myEntries.LoginPasswordNameEntryUpdate.Text, password,
				myEntries.LoginEntryUpdate.Text, myEntries.PasswordEntryUpdate.Text, accessToken)
			if err != nil {
				myLabels.LabelAlertLoginPasswordUpdate.SetText(errors.ErrLoginPasswordUpdate)
				b.log.Error(err)
			} else {
				dataTables.DataTblLoginPassword = table.UpdateRowLoginPassword(myEntries.LoginEntryUpdate.Text,
					myEntries.PasswordEntryUpdate.Text, dataTables.DataTblLoginPassword, myIndexes.IndexTblLoginPassword)
				b.log.Info("Логин-пароль изменен")

				myLabels.LabelAlertLoginPasswordUpdate.Hide()
				myForms.FormLoginPasswordUpdate.Refresh()
				b.window.SetContent(myTabs.ContainerTabs)
				b.window.Show()
			}
		} else {
			myLabels.LabelAlertLoginPasswordUpdate.SetText(errMsg)
			b.log.Error(errMsg)
		}
	})
}

func (b *Buttons) InitTextFormUpdateButton(myLabels *Labels, myEntries *Entries, myIndexes *Indexes,
	dataTables *DataTables, myForms *Forms, myTabs *Tabs, formValidator *FormValidator,
	accessToken model.Token, password string) {
	b.ButtonTextFormUpdate = widget.NewButton(labels.BtnUpdate, func() {
		myLabels.LabelAlertTextUpdate.Show()
		HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText,
			myLabels.LabelAlertCard, myLabels.LabelAlertBinary)

		errMsg, valid := formValidator.ValidateTextForm(Update)
		if valid {
			err := b.client.TextUpdate(myEntries.TextNameEntryUpdate.Text, password,
				myEntries.TextEntryUpdate.Text, accessToken)
			if err != nil {
				myLabels.LabelAlertTextUpdate.SetText(errors.ErrTextUpdate)
				b.log.Error(err)
			} else {
				dataTables.DataTblText = table.UpdateRowText(myEntries.TextEntryUpdate.Text, dataTables.DataTblText,
					myIndexes.IndexTblText)
				b.log.Info("Текст изменен")

				myLabels.LabelAlertTextUpdate.Hide()
				myForms.FormLoginPasswordUpdate.Refresh()
				b.window.SetContent(myTabs.ContainerTabs)
				b.window.Show()
			}
		} else {
			myLabels.LabelAlertTextUpdate.SetText(errMsg)
			b.log.Error(errMsg)
		}
	})
}

func (b *Buttons) InitCardFormUpdateButton(myLabels *Labels, myEntries *Entries, myIndexes *Indexes,
	dataTables *DataTables, myForms *Forms, myTabs *Tabs, formValidator *FormValidator,
	accessToken model.Token, password string) {
	b.ButtonCardFormUpdate = widget.NewButton(labels.BtnUpdate, func() {
		myLabels.LabelAlertCardUpdate.Show()
		HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText,
			myLabels.LabelAlertCard, myLabels.LabelAlertBinary)

		errMsg, valid := formValidator.ValidateCardForm(Update)
		if valid {
			err := b.client.CardUpdate(myEntries.CardNameEntryUpdate.Text, password,
				myEntries.PaymentSystemEntryUpdate.Text, myEntries.NumberEntryUpdate.Text,
				myEntries.HolderEntryUpdate.Text, myEntries.CvcEntryUpdate.Text, myEntries.EndDateEntryUpdate.Text, accessToken)
			if err != nil {
				myLabels.LabelAlertCardUpdate.SetText(errors.ErrCardUpdate)
				b.log.Error(err)
			} else {

				dataTables.DataTblCard = table.UpdateRowCard(myEntries.PaymentSystemEntryUpdate.Text,
					myEntries.NumberEntryUpdate.Text, myEntries.HolderEntryUpdate.Text,
					myEntries.CvcEntryUpdate.Text, myEntries.EndDateEntryUpdate.Text,
					dataTables.DataTblCard, myIndexes.IndexTblCard)
				b.log.Info("Карта изменена")

				myLabels.LabelAlertCardUpdate.Hide()
				myForms.FormCardUpdate.Refresh()
				b.window.SetContent(myTabs.ContainerTabs)
				b.window.Show()
			}
		} else {
			myLabels.LabelAlertCardUpdate.SetText(errMsg)
			b.log.Error(errMsg)
		}
	})
}

func (b *Buttons) InitBinaryUploadButton(myLabels *Labels, dataTables *DataTables, myTabs *Tabs,
	accessToken model.Token, password string) {
	b.ButtonBinaryUpload = widget.NewButton(labels.BtnUploadBinary, func() {
		fileDialog := dialog.NewFileOpen(
			func(r fyne.URIReadCloser, err error) {
				if r != nil {
					myLabels.LabelAlertBinary.Show()
					exist := table.SearchByColumn(dataTables.DataTblBinary, 0, r.URI().Name()) // search in map
					if exist {
						myLabels.LabelAlertBinary.SetText(errors.ErrBinaryExist)
						b.log.Error(myLabels.LabelAlertBinary.Text)
					}

					data, err := io.ReadAll(r)
					if len(data) > b.client.GetConfig().FileSize {
						myLabels.LabelAlertBinary.SetText(errors.ErrFileSize + "Размер загружаемого файла: " + strconv.Itoa(len(data)) + " байт")
						b.log.Error(errors.ErrFileSize + "Размер загружаемого файла: " + strconv.Itoa(len(data)) + " байт")
					} else if err != nil {
						myLabels.LabelAlertBinary.SetText(errors.ErrUpload)
						b.log.Error(err)
					} else {
						name, err := b.client.FileUpload(r.URI().Name(), password, data, accessToken)
						if err != nil {
							myLabels.LabelAlertBinary.SetText(errors.ErrUpload)
							b.log.Error(err)
						} else {
							dataTables.DataTblBinary = append(dataTables.DataTblBinary, []string{name, time.Now().Format(consts.DateAndTimeFormat)})
							b.log.Info("Файл добавлен: " + name)

							myLabels.LabelAlertBinary.Hide()
							b.window.SetContent(myTabs.ContainerTabs)
							b.window.Show()
						}
					}
				}
			}, b.window)
		fileDialog.Show()
	})
}

func (b *Buttons) InitBinaryDownloadButton(myLabels *Labels, myIndexes *Indexes,
	accessToken model.Token, password string) {
	b.ButtonBinaryDownload = widget.NewButton(labels.BtnDownloadBinary, func() {
		if myIndexes.IndexTblBinary > 0 {
			err := b.client.FileDownload(myIndexes.SelectedRowTblBinary[0], password, accessToken)
			if err != nil {
				myLabels.LabelAlertBinary.SetText(errors.ErrLogin)
				b.log.Error(err)
			}
			b.log.Info(myIndexes.IndexTblBinary)
			b.log.Info(myIndexes.SelectedRowTblBinary[0])
		} else {
			b.log.Error(errors.ErrBinaryTblIndexDownload)
			myLabels.LabelAlertBinary.Show()
			myLabels.LabelAlertBinary.SetText(errors.ErrBinaryTblIndexDownload)
		}
	})
}

func (b *Buttons) InitTopBackButton(myLabels *Labels, myEntries *Entries, myTabs *Tabs) {
	b.ButtonTopBack = widget.NewButton(labels.BtnBack, func() {
		ClearLoginPassword(myEntries.LoginPasswordNameEntryCreate, myEntries.LoginPasswordDescriptionEntryCreate,
			myEntries.LoginEntryCreate, myEntries.PasswordEntryCreate)
		ClearText(myEntries.TextNameEntryCreate, myEntries.TextDescriptionEntryCreate, myEntries.TextEntryCreate)
		ClearCard(myEntries.CardNameEntryCreate, myEntries.CardDescriptionEntryCreate, myEntries.PaymentSystemEntryCreate,
			myEntries.NumberEntryCreate, myEntries.HolderEntryCreate, myEntries.EndDateEntryCreate, myEntries.CvcEntryCreate)
		ClearLoginPassword(myEntries.LoginPasswordNameEntryUpdate, myEntries.LoginPasswordDescriptionEntryUpdate,
			myEntries.LoginEntryUpdate, myEntries.PasswordEntryUpdate)
		ClearText(myEntries.TextNameEntryUpdate, myEntries.TextDescriptionEntryUpdate, myEntries.TextEntryUpdate)
		ClearCard(myEntries.CardNameEntryUpdate, myEntries.CardDescriptionEntryUpdate, myEntries.PaymentSystemEntryUpdate,
			myEntries.NumberEntryUpdate, myEntries.HolderEntryUpdate, myEntries.EndDateEntryUpdate, myEntries.CvcEntryUpdate)
		myLabels.LabelAlertLoginPasswordCreate.Hide()
		myLabels.LabelAlertLoginPasswordUpdate.Hide()
		myLabels.LabelAlertTextCreate.Hide()
		myLabels.LabelAlertTextUpdate.Hide()
		myLabels.LabelAlertCardCreate.Hide()
		myLabels.LabelAlertCardUpdate.Hide()

		b.window.SetContent(myTabs.ContainerTabs)
		b.window.Resize(fyne.NewSize(windows.WindowMainWidth.Size(), windows.WindowMainHeight.Size()))
		b.window.Show()
	})
}

func (b *Buttons) InitAuthButton(myLabels *Labels, myEntries *Entries, myTabs *Tabs, dataTables *DataTables,
	radioAuth *widget.RadioGroup, formValidator *FormValidator, accessToken model.Token, password string) {
	b.ButtonAuth = widget.NewButton(labels.BtnSubmit, func() {
		myLabels.LabelAlertAuth.Show()
		if radioAuth.Selected == labels.RadioBtnLogin {
			errMsg, valid := formValidator.ValidateLoginForm()
			if valid {
				accessToken, err := b.client.Authentication(myEntries.UsernameLoginEntry.Text, myEntries.PasswordLoginEntry.Text)
				if err != nil {
					myLabels.LabelAlertAuth.SetText(errors.ErrLogin)
					b.log.Error(err)
				} else {
					password = myEntries.PasswordLoginEntry.Text
					dataTables.DataTblText, dataTables.DataTblCard, dataTables.DataTblLoginPassword, dataTables.DataTblBinary, err =
						b.client.Synchronization(password, accessToken)
					if err != nil {
						myLabels.LabelAlertAuth.SetText(errors.ErrLogin)
						b.log.Error(err)
					} else {
						b.window.SetContent(myTabs.ContainerTabs)
						b.window.Resize(fyne.NewSize(windows.WindowMainWidth.Size(), windows.WindowMainHeight.Size()))
						b.window.Show()
					}
				}
			} else {
				myLabels.LabelAlertAuth.SetText(errMsg)
				b.log.Error(errMsg)
			}
		}
		if radioAuth.Selected == labels.RadioBtnRegistration {
			errMsg, valid := formValidator.ValidateRegistrationForm()
			if valid {
				exist, err := b.client.UserExist(myEntries.UsernameRegistrationEntry.Text)
				if err != nil {
					myLabels.LabelAlertAuth.SetText(errors.ErrRegistration)
					b.log.Error(err)
				}
				if exist {
					myLabels.LabelAlertAuth.SetText(errors.ErrUserExist)
					b.log.Error(errors.ErrUserExist)
				} else {
					accessToken, err = b.client.Registration(myEntries.UsernameRegistrationEntry.Text, myEntries.PasswordRegistrationEntry.Text)
					if err != nil {
						myLabels.LabelAlertAuth.SetText(errors.ErrRegistration)
						b.log.Error(err)
					} else {
						password = myEntries.PasswordRegistrationEntry.Text
						b.window.SetContent(myTabs.ContainerTabs)
						b.window.Resize(fyne.NewSize(windows.WindowMainWidth.Size(), windows.WindowMainHeight.Size()))
						b.window.Show()
					}
				}
			} else {
				myLabels.LabelAlertAuth.SetText(errMsg)
				b.log.Error(errMsg)
			}
		}
	})
}

func (b *Buttons) InitLoginPasswordCreateButton(myLabels *Labels, myEntries *Entries, myForms *Forms, myTabs *Tabs,
	dataTables *DataTables, formValidator *FormValidator, accessToken model.Token, password string) {
	b.ButtonLoginPasswordCreate = widget.NewButton(labels.BtnAdd, func() {
		myLabels.LabelAlertLoginPasswordCreate.Show()
		HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText, myLabels.LabelAlertCard, myLabels.LabelAlertBinary)
		exist := table.SearchByColumn(dataTables.DataTblLoginPassword, 0, myEntries.LoginPasswordNameEntryCreate.Text) // search in map
		if exist {
			myLabels.LabelAlertLoginPassword.SetText(errors.ErrLoginPasswordExist)
			b.log.Error(myLabels.LabelAlertLoginPassword.Text)
		}
		errMsg, valid := formValidator.ValidateLoginPasswordForm(Create)
		if valid {
			err := b.client.LoginPasswordCreate(myEntries.LoginPasswordNameEntryCreate.Text,
				myEntries.LoginPasswordDescriptionEntryCreate.Text, password,
				myEntries.LoginEntryCreate.Text, myEntries.PasswordEntryCreate.Text, accessToken)
			if err != nil {
				myLabels.LabelAlertLoginPasswordCreate.SetText(errors.ErrLoginPasswordCreate)
				b.log.Error(err)
			} else {
				dataTables.DataTblLoginPassword = append(dataTables.DataTblLoginPassword,
					[]string{myEntries.LoginPasswordNameEntryCreate.Text, myEntries.LoginPasswordDescriptionEntryCreate.Text,
						myEntries.LoginEntryCreate.Text, myEntries.PasswordEntryCreate.Text, time.Now().Format(consts.DateAndTimeFormat),
						time.Now().Format(consts.DateAndTimeFormat)})

				ClearLoginPassword(myEntries.LoginPasswordNameEntryCreate, myEntries.LoginPasswordDescriptionEntryCreate,
					myEntries.LoginEntryCreate, myEntries.PasswordEntryCreate)
				b.log.Info("Логин-пароль добавлен")

				myLabels.LabelAlertLoginPasswordCreate.Hide()
				myForms.FormLoginPasswordCreate.Refresh()
				b.window.SetContent(myTabs.ContainerTabs)
				b.window.Show()
			}
		} else {
			myLabels.LabelAlertLoginPasswordCreate.SetText(errMsg)
			b.log.Error(errMsg)
		}
	})
}

func (b *Buttons) InitTextCreateButton(myLabels *Labels, myEntries *Entries, myForms *Forms, myTabs *Tabs,
	dataTables *DataTables, formValidator *FormValidator, accessToken model.Token, password string) {
	b.ButtonTextCreate = widget.NewButton(labels.BtnAdd, func() {
		myLabels.LabelAlertTextCreate.Show()
		HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText, myLabels.LabelAlertCard, myLabels.LabelAlertBinary)
		exist := table.SearchByColumn(dataTables.DataTblText, 0, myEntries.TextNameEntryCreate.Text) // search in map
		if exist {
			myLabels.LabelAlertText.SetText(errors.ErrTextExist)
			b.log.Error(myLabels.LabelAlertText)
		}
		errMsg, valid := formValidator.ValidateTextForm(Create)
		if valid {
			err := b.client.TextCreate(myEntries.TextNameEntryCreate.Text, myEntries.TextDescriptionEntryCreate.Text,
				password, myEntries.TextEntryCreate.Text, accessToken)
			if err != nil {
				myLabels.LabelAlertTextCreate.SetText(errors.ErrTextCreate)
				b.log.Error(err)
			} else {
				dataTables.DataTblText = append(dataTables.DataTblText, []string{myEntries.TextNameEntryCreate.Text,
					myEntries.TextDescriptionEntryCreate.Text, myEntries.TextEntryCreate.Text,
					time.Now().Format(consts.DateAndTimeFormat), time.Now().Format(consts.DateAndTimeFormat)})

				ClearText(myEntries.TextNameEntryCreate, myEntries.TextDescriptionEntryCreate, myEntries.TextEntryCreate)
				b.log.Info("Текст добавлен")

				myLabels.LabelAlertTextCreate.Hide()
				myForms.FormTextCreate.Refresh()
				b.window.SetContent(myTabs.ContainerTabs)
				b.window.Show()
			}
		} else {
			myLabels.LabelAlertTextCreate.SetText(errMsg)
			b.log.Error(errMsg)
		}
	})
}

func (b *Buttons) InitCardCreateButton(myLabels *Labels, myEntries *Entries, myForms *Forms, myTabs *Tabs,
	dataTables *DataTables, formValidator *FormValidator, accessToken model.Token, password string) {
	b.ButtonCardCreate = widget.NewButton(labels.BtnAdd, func() {
		myLabels.LabelAlertCardCreate.Show()
		HideLabelsTab(myLabels.LabelAlertLoginPassword, myLabels.LabelAlertText, myLabels.LabelAlertCard, myLabels.LabelAlertBinary)
		exist := table.SearchByColumn(dataTables.DataTblCard, 0, myEntries.CardNameEntryCreate.Text) // search in map
		if exist {
			myLabels.LabelAlertCard.SetText(errors.ErrCardExist)
			b.log.Print(myLabels.LabelAlertCard)
		}
		errMsg, valid := formValidator.ValidateCardForm(Create)
		if valid {
			err := b.client.CardCreate(myEntries.CardNameEntryCreate.Text, myEntries.CardDescriptionEntryCreate.Text, password,
				myEntries.PaymentSystemEntryCreate.Text, myEntries.NumberEntryCreate.Text, myEntries.HolderEntryCreate.Text,
				myEntries.CvcEntryCreate.Text, myEntries.EndDateEntryCreate.Text, accessToken)
			if err != nil {
				myLabels.LabelAlertCardCreate.SetText(errors.ErrCardCreate)
				b.log.Error(err)
			} else {
				dataTables.DataTblCard = append(dataTables.DataTblCard, []string{myEntries.CardNameEntryCreate.Text, myEntries.CardDescriptionEntryCreate.Text,
					myEntries.PaymentSystemEntryCreate.Text, myEntries.NumberEntryCreate.Text, myEntries.HolderEntryCreate.Text,
					myEntries.CvcEntryCreate.Text, myEntries.EndDateEntryCreate.Text, time.Now().Format(consts.DateAndTimeFormat),
					time.Now().Format(consts.DateAndTimeFormat)})

				ClearCard(myEntries.CardNameEntryCreate, myEntries.CardDescriptionEntryCreate,
					myEntries.PaymentSystemEntryCreate, myEntries.NumberEntryCreate, myEntries.HolderEntryCreate,
					myEntries.EndDateEntryCreate, myEntries.CvcEntryCreate)
				b.log.Info("Карта добавлена")

				myLabels.LabelAlertCardCreate.Hide()
				myForms.FormCardCreate.Refresh()
				b.window.SetContent(myTabs.ContainerTabs)
				b.window.Show()
			}
		} else {
			myLabels.LabelAlertCardCreate.SetText(errMsg)
			b.log.Error(errMsg)
		}
	})
}
