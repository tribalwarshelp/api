package utils

func SafeBoolPointer(value *bool, defaultValue bool) bool {
	if value == nil {
		return defaultValue
	}
	return *value
}

func SafeIntPointer(value *int, defaultValue int) int {
	if value == nil {
		return defaultValue
	}
	return *value
}

func SafeStrPointer(value *string, defaultValue string) string {
	if value == nil {
		return defaultValue
	}
	return *value
}
