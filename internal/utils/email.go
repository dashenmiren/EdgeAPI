package utils

import "regexp"

var emailReg = regexp.MustCompile(`(?i)^[a-z\d]+([._+-]*[a-z\d]+)*@([a-z\d]+[a-z\d-]*[a-z\d]+\.)+[a-z\d]+$`)

// ValidateEmail 校验电子邮箱格式
func ValidateEmail(email string) bool {
	return emailReg.MatchString(email)
}
