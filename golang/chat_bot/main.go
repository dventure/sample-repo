package main

import (
	"fmt"
	"strings"
)

func main() {

	for {
		fmt.Println(">You : ")
		var name string
		fmt.Scanln(&name)
		name = strings.ToLower(name)
		name = strings.TrimSpace(name)

		chatMap := map[string]string{
			"hi":    "Hi, How can I help you?",
			"hello": "Hello, How can I help you?",
			"name":  "I am bot.... Go-bot",
			"bye":   "Catch you later!",
			"quit":  "Bye, see you later!",
			"help":  "Allowed commands : \nHello  Name  Bye",
		}

		fmt.Println("Bot : ")
		if chatMap[name] != "" {
			fmt.Println(chatMap[name])
		} else {
			fmt.Println("Use 'help' command for allowed keywords")
		}
		if strings.Contains(name, "bye") || strings.Contains(name, "quit") {
			break
		}
	}
}
