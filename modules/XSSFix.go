package modules

import "strings"

func XSSFix(input string) string {
	a := strings.ReplaceAll(input, "<", "&lt;")
	a = strings.ReplaceAll(a, ">", "&gt;")
	a = strings.ReplaceAll(a, "(", "&#40;")
	a = strings.ReplaceAll(a, ")", "&#41;")
	a = strings.ReplaceAll(a, "\"", "&quot;")
	a = strings.ReplaceAll(a, "'", "&#x27;")
	a = strings.ReplaceAll(a, "/", "&#x2F;")

	a = strings.ReplaceAll(a, "\n", "<br>")
	return a
}
