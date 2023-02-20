package utils

import (
	"strconv"
	"strings"
)

func GetPathValue(obj interface{}, path string) interface{} {
	keys := strings.Split(path, ".")

	for _, key := range keys {
		switch val := obj.(type) {
		case map[string]interface{}:
			obj = val[key]
		case []interface{}:
			idx, err := strconv.Atoi(key)
			if err != nil {
				return nil
			}
			if idx < 0 || idx >= len(val) {
				return nil
			}
			obj = val[idx]
		default:
			return nil
		}
	}

	return obj
}
