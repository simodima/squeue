package sqs

import "os"

func CheckEnvVariables(vars ...string) bool {
	for _, v := range vars {
		if os.Getenv(v) == "" {
			return false
		}
	}

	return true
}
