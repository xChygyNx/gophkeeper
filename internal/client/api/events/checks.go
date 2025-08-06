package events

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type ValidateCardData struct {
	Cvc         string
	timeEndDate string
}

type ValidateResult struct {
	intCvc  int
	endDate time.Time
}

type CardChecker struct {
	logger     *logrus.Logger
	dateFormat string
	cardData   *ValidateCardData
	result     *ValidateResult
}

func NewCardChecker(logger *logrus.Logger, dateFormat string, cardData *ValidateCardData) *CardChecker {
	return &CardChecker{
		logger:     logger,
		dateFormat: dateFormat,
		cardData:   cardData,
		result:     &ValidateResult{},
	}
}

func (cc *CardChecker) RunChecks() error {
	checks := []func(data *ValidateCardData) error{cc.checkCVCCode, cc.checkTimeEndDate}

	for _, check := range checks {
		err := check(cc.cardData)
		if err != nil {
			funcName := runtime.FuncForPC(reflect.ValueOf(check).Pointer()).Name()
			return fmt.Errorf("error in check %s: %w", funcName, err)
		}
	}
	return nil
}

func (cc *CardChecker) checkCVCCode(data *ValidateCardData) error {
	intCvc, err := strconv.Atoi(data.Cvc)
	if err != nil {
		myErr := fmt.Errorf("invalid cvc code of card: %w", err)
		cc.logger.Error(myErr)
		return myErr
	}
	cc.result.intCvc = intCvc
	return nil
}

func (cc *CardChecker) checkTimeEndDate(data *ValidateCardData) error {
	timeEndDate, err := time.Parse(cc.dateFormat, data.timeEndDate)
	if err != nil {
		myErr := fmt.Errorf("invalid end time of card: %w", err)
		cc.logger.Error(myErr)
		return myErr
	}
	cc.result.endDate = timeEndDate
	return nil
}
