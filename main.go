package main

import (
	"NeuroNET/internal/neuronetserver"
	"NeuroNET/internal/neuronetserver/configs"
)

func main() {
	neuronetserver.New(configs.LocalConfigFile).Run()
}
