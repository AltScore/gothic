package xlogger

import "regexp"

func HideKey(key string) string {
	if key == "" {
		return "<NONE>"
	} else if len(key) < 2 {
		return "*"
	} else {
		return key[0:2] + "..." + key[len(key)-2:]
	}
}

// HideMongoSecret removes the password in the mongo URI
func HideMongoSecret(mongoUri string) string {
	return regexp.MustCompile(":([^@/]*@)+").ReplaceAllString(mongoUri, ":***@")
}
