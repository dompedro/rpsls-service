package main

import (
	base "rpsls"
	"rpsls/rpslsapi"
	"rpsls/rpslsapi/logger"
)

func main() {
	logger.SetUp()
	rpslsapi.LoadConfig()
	server, cleanup := base.InitServer()
	defer cleanup()
	server.Start()
}
