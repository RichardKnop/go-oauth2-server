package oauth

// Helpful function similar to "x in y" Python construct
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
