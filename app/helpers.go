package app

import (
	"github.com/Chatted-social/backend/storage"
	"regexp"
	"strconv"
)

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

func UsernameExists(db *storage.DB, u string) (bool, error) {

	if u == "" {
		return false, nil
	}

	user, err := db.Users.ExistsByUsername(u)

	if err != nil {
		return false, err
	}

	channel, err := db.Channels.ExistsByUsername(u)

	if err != nil {
		return false, err
	}

	return user || channel, nil

}

func UsernameValidator(u string) bool {
	var re = regexp.MustCompile(`^[a-zA-Z]+([_ -]?[a-zA-Z0-9])*$`)

	if len(re.FindStringIndex(u)) > 0 {
		return true
	}

	return false

}
