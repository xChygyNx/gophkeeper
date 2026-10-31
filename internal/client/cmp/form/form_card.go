package form

import (
	"fyne.io/fyne/v2/widget"

	"github.com/xChygyNx/gophkeeper/internal/client/storage/labels"
)

func GetFormCard(name *widget.Entry, cardDescription *widget.Entry, paymentSystem *widget.Entry,
	number *widget.Entry, holder *widget.Entry, endDate *widget.Entry, cvc *widget.Entry) *widget.Form {
	formCard := widget.NewForm(
		widget.NewFormItem(labels.NameItem, name),
		widget.NewFormItem(labels.DescriptionItem, cardDescription),
		widget.NewFormItem(labels.PaymentSystemItem, paymentSystem),
		widget.NewFormItem(labels.NumberItem, number),
		widget.NewFormItem(labels.HolderItem, holder),
		widget.NewFormItem(labels.EndDateItem, endDate),
		widget.NewFormItem(labels.CVCItem, cvc),
	)
	return formCard
}
