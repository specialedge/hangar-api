package events

import "fmt"

var DEBUG = true

func Info(key string, message string) {
	fmt.Println(key + " : " + message)
}

func Debug(key string, message string) {
	if DEBUG {
		fmt.Println(key + " : " + message)
	}
}
