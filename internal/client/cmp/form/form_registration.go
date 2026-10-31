package form

import (
	"fyne.io/fyne/v2/widget"

	"github.com/xChygyNx/gophkeeper/internal/client/storage/labels"
)

func GetFormRegistration(UsernameRegistration *widget.Entry, PasswordRegistration *widget.Entry, NewPasswordEntryRegistration *widget.Entry) *widget.Form {
	formRegistration := widget.NewForm(
		widget.NewFormItem(labels.UsernameItem, UsernameRegistration),
		widget.NewFormItem(labels.PasswordItem, PasswordRegistration),
		widget.NewFormItem(labels.ConfirmPasswordItem, NewPasswordEntryRegistration),
	)
	return formRegistration
}
