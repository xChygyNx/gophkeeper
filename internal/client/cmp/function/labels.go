package function

import "fyne.io/fyne/v2/widget"

type Labels struct {
	LabelAlertAuth                *widget.Label
	LabelAlertLoginPassword       *widget.Label
	LabelAlertLoginPasswordCreate *widget.Label
	LabelAlertLoginPasswordUpdate *widget.Label
	LabelAlertText                *widget.Label
	LabelAlertTextCreate          *widget.Label
	LabelAlertTextUpdate          *widget.Label
	LabelAlertCard                *widget.Label
	LabelAlertCardCreate          *widget.Label
	LabelAlertCardUpdate          *widget.Label
	LabelAlertBinary              *widget.Label
}

func InitLabels() *Labels {
	labels := &Labels{
		LabelAlertAuth:                widget.NewLabel(""),
		LabelAlertLoginPassword:       widget.NewLabel(""),
		LabelAlertLoginPasswordCreate: widget.NewLabel(""),
		LabelAlertLoginPasswordUpdate: widget.NewLabel(""),
		LabelAlertText:                widget.NewLabel(""),
		LabelAlertTextCreate:          widget.NewLabel(""),
		LabelAlertTextUpdate:          widget.NewLabel(""),
		LabelAlertCard:                widget.NewLabel(""),
		LabelAlertCardCreate:          widget.NewLabel(""),
		LabelAlertCardUpdate:          widget.NewLabel(""),
		LabelAlertBinary:              widget.NewLabel(""),
	}

	labels.LabelAlertAuth.Hide()
	labels.LabelAlertLoginPassword.Hide()
	labels.LabelAlertLoginPasswordCreate.Hide()
	labels.LabelAlertLoginPasswordUpdate.Hide()
	labels.LabelAlertText.Hide()
	labels.LabelAlertTextCreate.Hide()
	labels.LabelAlertTextUpdate.Hide()
	labels.LabelAlertCard.Hide()
	labels.LabelAlertCardCreate.Hide()
	labels.LabelAlertCardUpdate.Hide()
	labels.LabelAlertBinary.Hide()

	return labels
}
