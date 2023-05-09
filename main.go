package main

import (
	"github.com/Peterpig/mini_godis/config"
	"github.com/Peterpig/mini_godis/lib/files"
	"github.com/Peterpig/mini_godis/lib/logger"
	"github.com/Peterpig/mini_godis/lib/utils"
)

var banner = `
   ______          ___
  / ____/___  ____/ (_)____
 / / __/ __ \/ __  / / ___/
/ /_/ / /_/ / /_/ / (__  )
\____/\____/\__,_/_/____/
`

var defultProperties = &config.ServerProperties{
	Bind:       "0.0.0.0",
	Port:       6399,
	AppendOnly: false,
	MaxClients: 1000,
	RunID:      utils.RandString(40),
}

func main() {
	print(banner)
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if files.FileExists("redis.conf") {
		config.SetupConfig("redis.conf")
	} else {
		config.Properties = defultProperties
	}
}
