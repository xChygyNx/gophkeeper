package function

import "fyne.io/fyne/v2/widget"

func ClearText(textNameEntry *widget.Entry, textEntry *widget.Entry, textDescriptionEntry *widget.Entry) {
	textNameEntry.Text = ""
	textEntry.Text = ""
	textDescriptionEntry.Text = ""
}

func ClearCard(cardNameEntry *widget.Entry, cardDescriptionEntry *widget.Entry, paymentSystemEntry *widget.Entry,
	numberEntry *widget.Entry, holderEntry *widget.Entry, endDateEntry *widget.Entry, cvcEntry *widget.Entry) {
	cardNameEntry.SetText("")
	cardDescriptionEntry.SetText("")
	cardNameEntry.Text = ""
	paymentSystemEntry.Text = ""
	numberEntry.Text = ""
	holderEntry.Text = ""
	endDateEntry.Text = ""
	cvcEntry.Text = ""
}

func ClearLoginPassword(loginPasswordNameEntry *widget.Entry, loginPasswordDescriptionEntry *widget.Entry,
	loginEntry *widget.Entry, passwordEntry *widget.Entry) {
	loginPasswordNameEntry.SetText("")
	loginPasswordDescriptionEntry.SetText("")
	loginEntry.SetText("")
	passwordEntry.SetText("")
}
