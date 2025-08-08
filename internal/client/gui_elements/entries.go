package gui_elements

import "fyne.io/fyne/v2/widget"

type Entries struct {
	UsernameLoginEntry                    *widget.Entry
	PasswordLoginEntry                    *widget.Entry
	UsernameRegistrationEntry             *widget.Entry
	PasswordRegistrationEntry             *widget.Entry
	PasswordConfirmationRegistrationEntry *widget.Entry
	TextNameEntryCreate                   *widget.Entry
	TextDescriptionEntryCreate            *widget.Entry
	TextEntryCreate                       *widget.Entry
	TextNameEntryUpdate                   *widget.Entry
	TextDescriptionEntryUpdate            *widget.Entry
	TextEntryUpdate                       *widget.Entry
	LoginPasswordNameEntryCreate          *widget.Entry
	LoginPasswordDescriptionEntryCreate   *widget.Entry
	LoginEntryCreate                      *widget.Entry
	PasswordEntryCreate                   *widget.Entry
	LoginPasswordNameEntryUpdate          *widget.Entry
	LoginPasswordDescriptionEntryUpdate   *widget.Entry
	LoginEntryUpdate                      *widget.Entry
	PasswordEntryUpdate                   *widget.Entry
	CardNameEntryCreate                   *widget.Entry
	CardDescriptionEntryCreate            *widget.Entry
	PaymentSystemEntryCreate              *widget.Entry
	NumberEntryCreate                     *widget.Entry
	HolderEntryCreate                     *widget.Entry
	EndDateEntryCreate                    *widget.Entry
	CvcEntryCreate                        *widget.Entry
	CardNameEntryUpdate                   *widget.Entry
	CardDescriptionEntryUpdate            *widget.Entry
	PaymentSystemEntryUpdate              *widget.Entry
	NumberEntryUpdate                     *widget.Entry
	HolderEntryUpdate                     *widget.Entry
	EndDateEntryUpdate                    *widget.Entry
	CvcEntryUpdate                        *widget.Entry
}

func InitEntries() *Entries {
	result := &Entries{
		UsernameLoginEntry:                    widget.NewEntry(),
		PasswordLoginEntry:                    widget.NewPasswordEntry(),
		UsernameRegistrationEntry:             widget.NewEntry(),
		PasswordRegistrationEntry:             widget.NewPasswordEntry(),
		PasswordConfirmationRegistrationEntry: widget.NewPasswordEntry(),
		TextNameEntryCreate:                   widget.NewEntry(),
		TextDescriptionEntryCreate:            widget.NewEntry(),
		TextEntryCreate:                       widget.NewEntry(),
		TextNameEntryUpdate:                   widget.NewEntry(),
		TextDescriptionEntryUpdate:            widget.NewEntry(),
		TextEntryUpdate:                       widget.NewEntry(),
		LoginPasswordNameEntryCreate:          widget.NewEntry(),
		LoginPasswordDescriptionEntryCreate:   widget.NewEntry(),
		LoginEntryCreate:                      widget.NewEntry(),
		PasswordEntryCreate:                   widget.NewEntry(),
		LoginPasswordNameEntryUpdate:          widget.NewEntry(),
		LoginPasswordDescriptionEntryUpdate:   widget.NewEntry(),
		LoginEntryUpdate:                      widget.NewEntry(),
		PasswordEntryUpdate:                   widget.NewEntry(),
		CardNameEntryCreate:                   widget.NewEntry(),
		CardDescriptionEntryCreate:            widget.NewEntry(),
		PaymentSystemEntryCreate:              widget.NewEntry(),
		NumberEntryCreate:                     widget.NewEntry(),
		HolderEntryCreate:                     widget.NewEntry(),
		EndDateEntryCreate:                    widget.NewEntry(),
		CvcEntryCreate:                        widget.NewEntry(),
		CardNameEntryUpdate:                   widget.NewEntry(),
		CardDescriptionEntryUpdate:            widget.NewEntry(),
		PaymentSystemEntryUpdate:              widget.NewEntry(),
		NumberEntryUpdate:                     widget.NewEntry(),
		HolderEntryUpdate:                     widget.NewEntry(),
		EndDateEntryUpdate:                    widget.NewEntry(),
		CvcEntryUpdate:                        widget.NewEntry(),
	}

	result.TextNameEntryUpdate.Disable()
	result.TextDescriptionEntryUpdate.Disable()
	result.LoginPasswordNameEntryUpdate.Disable()
	result.LoginPasswordDescriptionEntryUpdate.Disable()
	result.CardNameEntryUpdate.Disable()
	result.CardDescriptionEntryUpdate.Disable()

	return result
}
