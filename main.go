package main

import (
	"github.com/GitGert/Pipedrive-Devops-challenge/src/server"
	"github.com/GitGert/Pipedrive-Devops-challenge/src/utils"
)

func main() {
	utils.LoadEnvFile(".env")
	server.InitServer()
}
