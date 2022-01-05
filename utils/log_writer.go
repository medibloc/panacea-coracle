package utils

import "log"

func WriteLogger(_ int, err error) {
	if err != nil {
		log.Println("Write failed %v", err)
	}
}
