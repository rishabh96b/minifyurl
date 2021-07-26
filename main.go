package main

import (
	"github.com/rishabh96b/minifyurl/app"
	"github.com/rishabh96b/minifyurl/config"
)

func main() {
	config := config.GetMySQlDBConfig()

	app := &app.App{}
	app.Initialize(config)
	defer app.DB.Close()

	app.Run(":8000")
}
