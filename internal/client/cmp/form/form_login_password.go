package form

import (
	"fyne.io/fyne/v2/widget"

	"github.com/xChygyNx/gophkeeper/internal/client/storage/labels"
)

func GetFormLoginPassword(loginPasswordName *widget.Entry, loginPasswordDescription *widget.Entry, login *widget.Entry, password *widget.Entry) *widget.Form {
	formText := widget.NewForm(
		widget.NewFormItem(labels.NameItem, loginPasswordName),
		widget.NewFormItem(labels.DescriptionItem, loginPasswordDescription),
		widget.NewFormItem(labels.LoginItem, login),
		widget.NewFormItem(labels.PasswordItem, password),
	)
	return formText
}
