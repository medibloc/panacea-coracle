package testutil

import "strings"

func combinesKey(key ...string) string {
	return strings.Join(key, "|")
}
