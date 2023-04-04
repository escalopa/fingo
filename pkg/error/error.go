package pkgerror

import "log"

func CheckError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
