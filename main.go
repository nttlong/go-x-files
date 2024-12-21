package main

import (
	"fmt"

	application "github.com/unvs/libs/application"
)

func main() {

	app := application.AppContext

	app_path := app.AppPath

	fmt.Println(app_path)
	fmt.Println(app.Config)

}
