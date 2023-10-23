package validators

import "regexp"

func IsEmailValid(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return isValidFormat(email)
}

func isValidFormat(email string) bool {
	// Regular expression for email validation
	emailPattern := "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	rxEmail := regexp.MustCompile(emailPattern)
	return rxEmail.MatchString(email)
}
