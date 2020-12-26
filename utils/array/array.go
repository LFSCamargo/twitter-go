package array

// RemoveItem - Finds and remove an item from a array of strings
func RemoveItem(slice []string, item string) []string {
	for i, value := range slice {
		if value == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// FindItem - checks if that item exists inside the array
func FindItem(slice []string, val string) (string, bool) {
	for _, item := range slice {
		if item == val {
			return item, true
		}
	}
	return "", false
}
