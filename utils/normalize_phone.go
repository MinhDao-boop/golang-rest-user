package utils

import (
	"errors"
	"regexp"
	"strings"
)

var e164Regex = regexp.MustCompile("^\\+[1-9]\\d{7,14}$")

func isValidE164(phone string) bool {
	return e164Regex.MatchString(phone)
}

func NormalizePhone(phoneRaw string) (string, error) {
	phone := strings.TrimSpace(phoneRaw)

	regex := regexp.MustCompile(`[^\d+]`)
	phone = regex.ReplaceAllString(phone, "")

	if strings.HasPrefix(phone, "+") {
		if !isValidE164(phone) {
			return "", errors.New("invalid international phone number")
		}
		return phone, nil
	}

	if strings.HasPrefix(phone, "0") {
		phone = "+84" + phone[1:]
		if !isValidE164(phone) {
			return "", errors.New("invalid Vietnamese phone number")
		}
		return phone, nil
	}
	//
	//if strings.HasPrefix(phone, "84") {
	//	phone = "+" + phone
	//	if !isValidE164(phone) {
	//		return "", errors.New("invalid phone number")
	//	}
	//	return phone, nil
	//}
	return "", errors.New("unsupported phone format")
}
