package form

import (
	"fyne.io/fyne/v2/widget"

	"github.com/xChygyNx/gophkeeper/internal/client/storage/labels"
)

func GetFormLogin(username *widget.Entry, password *widget.Entry) *widget.Form {
	formLogin := widget.NewForm(
		widget.NewFormItem(labels.UsernameItem, username),
		widget.NewFormItem(labels.PasswordItem, password),
	)
	return formLogin
}
