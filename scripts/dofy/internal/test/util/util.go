package util

import "os"

func IsCI() bool {
	return os.Getenv("GITHUB_ACIONS") == "true"
}
