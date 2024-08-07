package utils


func HideEmail(email string) string {
	runeEmail := []rune(email)
	return string(runeEmail[:2]) + "**********@*******"
}

func HidePhone(phone string) string {
	runePhone := []rune(phone)
	return "********" + string(runePhone[len(runePhone)-2:])
}
