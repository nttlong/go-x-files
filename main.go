package main

import (
	"fmt"

	application "github.com/unvs/libs/application"
	looger "github.com/unvs/libs/loggers"
)

func main() {
	defer looger.HandlePanic()
	application.InitGlobalContext()

	app := application.AppContext

	app_path := app.AppPath
	key := app.Cacher.GetKey("app_path")
	println(key)
	app.Cacher.SetString("app_path", app_path, 0)
	test := app.Cacher.GetString("app_path")

	fmt.Println(test)

}
