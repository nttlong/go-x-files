// this file decribes the Application struct and Run() function
package application

import (
	"path/filepath"

	cacher "github.com/unvs/libs/cacher"

	config "github.com/unvs/libs/configReader"
)

// app struct
type Application struct {
	// app name
	Name string
	// version of the app will be set by the build process or env variable
	// in docker container
	Version string
	AppPath string
	Config  *config.Config
	Cacher  cacher.Cacher
}

func (app *Application) Init() {
	app.Name = "go-x-files"
	app.Version = "0.0.1"
	app.AppPath = config.GetAppPath()
	// move to parent directory to load config file
	app.AppPath = filepath.Dir(filepath.Dir(app.AppPath))

	appConfig := config.LoadConfig(app.AppPath + "/config.yml")

	app.Config = appConfig
	//cast mc to cacher interface

	// app.Cacher = &memcacher.MemcacheCacher{
	// 	Server: app.Config.CacheServer,
	// 	Prefix: app.Config.CachePrefix,
	// 	// set default expiry is 4 hours
	// 	Expiry: 14400,
	// }

}

var AppContext Application

func InitGlobalContext() {
	AppContext = Application{}
	AppContext.Init()

}
