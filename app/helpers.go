package app

func Ok() map[string]bool {
	return map[string]bool{"ok": true}
}

func Err(msg string) map[string]string {
	return map[string]string{"error": msg}
}