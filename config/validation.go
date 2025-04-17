package config

func ValidUsername(s string) bool {
	if len(s) < 8 || len(s) > 50 {
		return false
	}
	if !UsernameExp.MatchString(s) {
		return false
	}
	return true
}

func ValidEmail(s string) bool {
	if len(s) < 8 || len(s) > 50 {
		return false
	}
	if !EmailExp.MatchString(s) {
		return false
	}
	return true
}
