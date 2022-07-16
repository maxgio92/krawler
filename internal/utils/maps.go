package utils

func MergeMapsAndDeleteKeys(m1 map[string]interface{}, m2 map[string]interface{}, keysToDelete ...string) map[string]interface{} {
	result := MergeMaps(m1, m2)

	for _, k := range keysToDelete {
		delete(result, k)
	}

	return result
}

func MergeMaps(m1 map[string]interface{}, m2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {
		result[k] = v
	}

	return result
}
