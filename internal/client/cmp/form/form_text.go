package form

import (
	"fyne.io/fyne/v2/widget"

	"github.com/xChygyNx/gophkeeper/internal/client/storage/labels"
)

func GetFormText(textName *widget.Entry, textDescription *widget.Entry, text *widget.Entry) *widget.Form {
	formText := widget.NewForm(
		widget.NewFormItem(labels.NameItem, textName),
		widget.NewFormItem(labels.DescriptionItem, textDescription),
		widget.NewFormItem(labels.DataItem, text),
	)
	return formText
}
