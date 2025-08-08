package function

import (
	"fyne.io/fyne/v2/widget"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/xChygyNx/gophkeeper/internal/client/consts"
	"github.com/xChygyNx/gophkeeper/internal/client/gui_elements"
	"github.com/xChygyNx/gophkeeper/internal/client/service/algorithm"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	"github.com/xChygyNx/gophkeeper/internal/client/storage/errors"
)

type Mode string

const (
	Create Mode = "Create"
	Update Mode = "Update"
)

const userNameMinLength = 6

type FormValidator struct {
	*gui_elements.Entries
	emptyValue string
}

func NewFormValidator(entries *gui_elements.Entries) *FormValidator {
	return &FormValidator{
		Entries:    entries,
		emptyValue: "",
	}
}

func (fv *FormValidator) ValidateLoginForm() (string, bool) {
	if utf8.RuneCountInString(fv.UsernameLoginEntry.Text) < userNameMinLength {
		return errors.ErrUsernameIncorrect, false
	}
	if !encryption.VerifyPassword(fv.PasswordLoginEntry.Text) {
		return errors.ErrPasswordIncorrect, false
	}
	return fv.emptyValue, true
}

func (fv *FormValidator) ValidateRegistrationForm() (string, bool) {
	if utf8.RuneCountInString(fv.UsernameRegistrationEntry.Text) < userNameMinLength {
		return errors.ErrUsernameIncorrect, false
	}
	if !encryption.VerifyPassword(fv.PasswordRegistrationEntry.Text) {
		return errors.ErrPasswordIncorrect, false
	}
	if fv.PasswordRegistrationEntry.Text != fv.PasswordConfirmationRegistrationEntry.Text {
		return errors.ErrPasswordDifferent, false
	}
	return fv.emptyValue, true
}

func (fv *FormValidator) ValidateLoginPasswordForm(mode Mode) (string, bool) {
	var loginPasswordNameEntry *widget.Entry
	var loginPasswordDescriptionEntry *widget.Entry
	var loginEntry *widget.Entry
	var passwordEntry *widget.Entry

	if mode == Create {
		loginPasswordNameEntry = fv.LoginPasswordNameEntryCreate
		loginPasswordDescriptionEntry = fv.LoginPasswordDescriptionEntryCreate
		loginEntry = fv.LoginEntryCreate
		passwordEntry = fv.PasswordEntryCreate
	} else if mode == Update {
		loginPasswordNameEntry = fv.LoginPasswordNameEntryUpdate
		loginPasswordDescriptionEntry = fv.LoginPasswordDescriptionEntryUpdate
		loginEntry = fv.LoginEntryUpdate
		passwordEntry = fv.PasswordEntryUpdate
	}

	if loginPasswordNameEntry.Text == fv.emptyValue {
		return errors.ErrNameEmpty, false
	}
	if loginPasswordDescriptionEntry.Text == fv.emptyValue {
		return errors.ErrDescriptionEmpty, false
	}
	if loginEntry.Text == fv.emptyValue {
		return errors.ErrLoginEmpty, false
	}
	if passwordEntry.Text == fv.emptyValue {
		return errors.ErrPasswordEmpty, false
	}
	return fv.emptyValue, true
}

func (fv *FormValidator) ValidateTextForm(mode Mode) (string, bool) {
	var textNameEntry *widget.Entry
	var textDescriptionEntry *widget.Entry
	var textEntry *widget.Entry

	if mode == Create {
		textNameEntry = fv.TextNameEntryCreate
		textDescriptionEntry = fv.TextDescriptionEntryCreate
		textEntry = fv.TextEntryCreate
	} else if mode == Update {
		textNameEntry = fv.TextNameEntryUpdate
		textDescriptionEntry = fv.TextDescriptionEntryUpdate
		textEntry = fv.TextEntryUpdate
	}

	if textNameEntry.Text == fv.emptyValue {
		return errors.ErrNameEmpty, false
	}
	if textDescriptionEntry.Text == fv.emptyValue {
		return errors.ErrDescriptionEmpty, false
	}
	if textEntry.Text == fv.emptyValue {
		return errors.ErrTextEmpty, false
	}
	return fv.emptyValue, true
}

func (fv *FormValidator) ValidateCardForm(mode Mode) (string, bool) {
	var err error
	var cardNameEntry *widget.Entry
	var cardDescriptionEntry *widget.Entry
	var paymentSystemEntry *widget.Entry
	var numberEntry *widget.Entry
	var holderEntry *widget.Entry
	var endDateEntry *widget.Entry
	var cvcEntry *widget.Entry

	if mode == Create {
		cardNameEntry = fv.CardNameEntryCreate
		cardDescriptionEntry = fv.CardDescriptionEntryCreate
		paymentSystemEntry = fv.PaymentSystemEntryCreate
		numberEntry = fv.NumberEntryCreate
		holderEntry = fv.HolderEntryCreate
		endDateEntry = fv.EndDateEntryCreate
		cvcEntry = fv.CvcEntryCreate
	} else if mode == Update {
		cardNameEntry = fv.CardNameEntryUpdate
		cardDescriptionEntry = fv.CardDescriptionEntryUpdate
		paymentSystemEntry = fv.PaymentSystemEntryUpdate
		numberEntry = fv.NumberEntryUpdate
		holderEntry = fv.HolderEntryUpdate
		endDateEntry = fv.EndDateEntryUpdate
		cvcEntry = fv.CvcEntryUpdate
	}
	if cardNameEntry.Text == fv.emptyValue {
		return errors.ErrNameEmpty, false
	}
	if cardDescriptionEntry.Text == fv.emptyValue {
		return errors.ErrDescriptionEmpty, false
	}
	if paymentSystemEntry.Text == fv.emptyValue {
		return errors.ErrPaymentSystemEmpty, false
	}
	if numberEntry.Text == fv.emptyValue {
		return errors.ErrNumberEmpty, false
	}
	intNumber, err := strconv.Atoi(numberEntry.Text)
	if err != nil {
		return errors.ErrNumberIncorrect, false
	}
	if !algorithm.ValidLuhn(intNumber) {
		return errors.ErrNumberIncorrect, false
	}
	if holderEntry.Text == fv.emptyValue {
		return errors.ErrHolderEmpty, false
	}
	if endDateEntry.Text == fv.emptyValue {
		return errors.ErrEndDateEmpty, false
	} else {
		_, err = time.Parse(consts.DateFormat, endDateEntry.Text)
		if err != nil {
			return errors.ErrEndDateIncorrect, false
		}
	}
	if cvcEntry.Text == fv.emptyValue {
		return errors.ErrCvcEmpty, false
	} else {
		_, err = strconv.Atoi(cvcEntry.Text)
		if err != nil {
			return errors.ErrCvcIncorrect, false
		}
	}
	return fv.emptyValue, true
}
