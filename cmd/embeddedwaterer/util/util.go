package util

// MaybeLeadingComma returns a leading comma for every idx after 0.
// Used as a leading comma for JSON to add commas between items
func MaybeLeadingComma(idx int) string {
	if idx > 0 {
		return ","
	}
	return ""
}
