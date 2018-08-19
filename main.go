package main

import (
	"fmt"
	"friend-management/controller"
)

var appName = "friendManagement"

func main() {
	fmt.Printf("Starting %v\n", appName)
	controller.StartWebServer("8000")
}
