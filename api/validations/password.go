package validations

func ValidatePassword(password string) bool {
	return len(password) < 8
}
