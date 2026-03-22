package fsc

import (
	"strconv"
)

func firstString(m map[string]any, keys ...string) string {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		switch vv := v.(type) {
		case string:
			if vv != "" {
				return vv
			}
		case float64:
			return strconv.FormatInt(int64(vv), 10)
		case int:
			return strconv.Itoa(vv)
		}
	}
	return ""
}

func firstInt(m map[string]any, keys ...string) int {
	for _, key := range keys {
		v, ok := m[key]
		if !ok || v == nil {
			continue
		}
		switch vv := v.(type) {
		case float64:
			return int(vv)
		case int:
			return vv
		case string:
			parsed, err := strconv.Atoi(vv)
			if err == nil {
				return parsed
			}
		}
	}
	return 0
}
