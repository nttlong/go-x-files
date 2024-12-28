// this file decribes the Application struct and Run() function
package application

import (
	"context"

	cacher "github.com/unvs/libs/cacher"

	config "github.com/unvs/libs/configReader"
	dbx "github.com/unvs/libs/db/ctx"
)

// app struct
type Application struct {
	// app name
	Name string
	// version of the app will be set by the build process or env variable
	// in docker container
	Version    string
	AppPath    string
	Config     *config.Config
	Cacher     *cacher.Cacher
	AppContext *context.Context
	DB         *dbx.DBContext
}

func CreateApp() Application {
	app := Application{}
	return app
}
func (app *Application) SetAppPath(path string) Application {
	app.AppPath = path
	return *app
}
func (app *Application) LoadConfig() Application {
	if app.AppPath == "" {
		panic("AppPath is not set")
	}
	if app.AppContext == nil {
		panic("AppContext is not set")
	}
	app.Config = config.LoadConfig(app.AppPath + "/config.yml")

	return *app
}
func (app *Application) SetCacher(cacher cacher.Cacher) Application {
	app.Cacher = &cacher
	return *app
}
func (app *Application) SetContext(ctx context.Context) Application {
	app.AppContext = &ctx
	return *app
}
