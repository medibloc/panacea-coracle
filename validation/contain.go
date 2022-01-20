package validation

func Contains(strList []string, str string) bool {
	for _, v := range strList {
		if v == str {
			return true
		}
	}

	return false
}
