package gui_elements

import (
	"fyne.io/fyne/v2/widget"
	"github.com/xChygyNx/gophkeeper/internal/client/cmp/form"
)

type Forms struct {
	FormLogin               *widget.Form
	FormRegistration        *widget.Form
	FormLoginPasswordCreate *widget.Form
	FormLoginPasswordUpdate *widget.Form
	FormTextCreate          *widget.Form
	FormTextUpdate          *widget.Form
	FormCardCreate          *widget.Form
	FormCardUpdate          *widget.Form
}

func InitForms(entries *Entries) *Forms {
	return &Forms{
		FormLogin:        form.GetFormLogin(entries.UsernameLoginEntry, entries.PasswordLoginEntry),
		FormRegistration: form.GetFormRegistration(entries.UsernameRegistrationEntry, entries.PasswordRegistrationEntry, entries.PasswordConfirmationRegistrationEntry),

		FormLoginPasswordCreate: form.GetFormLoginPassword(entries.LoginPasswordNameEntryCreate, entries.LoginPasswordDescriptionEntryCreate,
			entries.LoginEntryCreate, entries.PasswordEntryCreate),
		FormLoginPasswordUpdate: form.GetFormLoginPassword(entries.LoginPasswordNameEntryUpdate, entries.LoginPasswordDescriptionEntryUpdate,
			entries.LoginEntryUpdate, entries.PasswordEntryUpdate),

		FormTextCreate: form.GetFormText(entries.TextNameEntryCreate, entries.TextDescriptionEntryCreate, entries.TextEntryCreate),
		FormTextUpdate: form.GetFormText(entries.TextNameEntryUpdate, entries.TextDescriptionEntryUpdate, entries.TextEntryUpdate),

		FormCardCreate: form.GetFormCard(entries.CardNameEntryCreate, entries.CardDescriptionEntryCreate, entries.PaymentSystemEntryCreate,
			entries.NumberEntryCreate, entries.HolderEntryCreate, entries.EndDateEntryCreate, entries.CvcEntryCreate),
		FormCardUpdate: form.GetFormCard(entries.CardNameEntryUpdate, entries.CardDescriptionEntryUpdate, entries.PaymentSystemEntryUpdate,
			entries.NumberEntryUpdate, entries.HolderEntryUpdate, entries.EndDateEntryUpdate, entries.CvcEntryUpdate),
	}
}
