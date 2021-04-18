package app

import "strconv"

func Ok() map[string]bool {
	return map[string]bool{"ok": true}
}

func Err(msg string) map[string]string {
	return map[string]string{"error": msg}
}

func StringSliceToInt(slice []string) (r []int) {
	for _, s := range slice {
		if i, err := strconv.Atoi(s); err == nil {
			r = append(r, i)
		}
	}
	return
}
