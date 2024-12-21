package main

import (
	"fmt"

	application "github.com/unvs/libs/application"
	cacher "github.com/unvs/libs/cacher"
	memmcacher "github.com/unvs/libs/cacher/memcacher"
)

func main() {

	app := application.AppContext
	app.RegisterService((*cacher.Cacher)(nil), memmcacher.NewMemMCacher)
	app.GetService(cacher.Cacher)
	app_path := app.AppPath

	fmt.Println(app_path)
	fmt.Println(app.Config)

}
