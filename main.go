package main

import (
	"neuronet/internal/neuronetserver"
	"neuronet/internal/neuronetserver/configs"
)

func main() {
	neuronetserver.New(configs.LocalConfigFile).Run()
}
