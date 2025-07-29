package function

import "fyne.io/fyne/v2/widget"

func SetLoginPasswordData(loginPasswordData []string, loginPasswordNameEntry *widget.Entry, loginPasswordDescriptionEntry *widget.Entry, loginEntry *widget.Entry, passwordEntry *widget.Entry) {
	loginPasswordNameEntry.SetText(loginPasswordData[0])
	loginPasswordDescriptionEntry.SetText(loginPasswordData[1])
	loginEntry.SetText(loginPasswordData[2])
	passwordEntry.SetText(loginPasswordData[3])
}

func SetTextData(textData []string, textNameEntry *widget.Entry, textDescriptionEntry *widget.Entry, textEntry *widget.Entry) {
	textNameEntry.SetText(textData[0])
	textDescriptionEntry.SetText(textData[1])
	textEntry.SetText(textData[2])
}

func SetCardData(card []string, cardNameEntry *widget.Entry, cardDescriptionEntry *widget.Entry, paymentSystemEntry *widget.Entry, numberEntry *widget.Entry, holderEntry *widget.Entry, endDateEntry *widget.Entry, cvcEntry *widget.Entry) {
	cardNameEntry.SetText(card[0])
	cardDescriptionEntry.SetText(card[1])
	paymentSystemEntry.SetText(card[2])
	numberEntry.SetText(card[3])
	holderEntry.SetText(card[4])
	endDateEntry.SetText(card[5])
	cvcEntry.SetText(card[6])
}
