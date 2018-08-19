package main

import (
	"fmt"
	"friend-management/controller"
	"friend-management/repository"
)

var appName = "friendManagement"

func main() {
	fmt.Printf("Starting %v\n", appName)
	initializeDB()
	controller.StartWebServer("8080")
}

// Creates instance and calls the OpenBoltDb and Seed funcs
func initializeDB() {
	controller.UserRepo = &repository.Repository{}
	controller.UserRepo.OpenBoltDb()
	controller.UserRepo.Seed()
}
