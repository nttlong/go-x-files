// this file decribes the Application struct and Run() function
package application

import (
	"path/filepath"
	"reflect"

	config "github.com/unvs/libs/configReader"
)

// app struct
type Application struct {
	// app name
	Name string
	// version of the app will be set by the build process or env variable
	// in docker container
	Version  string
	AppPath  string
	Config   *config.Config
	services map[reflect.Type]interface{}
}

func (app *Application) Init() (err error) {
	app.Name = "go-x-files"
	app.Version = "0.0.1"
	app.AppPath = config.GetAppPath()
	// move to parent directory to load config file
	app.AppPath = filepath.Dir(filepath.Dir(app.AppPath))

	appConfig, err := config.LoadConfig(app.AppPath + "/config.yml")
	if err != nil {
		return err
	}
	app.Config = appConfig
	return nil

}

var AppContext Application

func createInstance[T any]() T {
	var zero T
	if reflect.TypeOf(zero).Kind() == reflect.Interface {
		return zero // Return nil if it is interface
	}
	instance := reflect.New(reflect.TypeOf(zero)).Elem().Interface().(T)
	return instance
}

func init() {
	AppContext = Application{}
	err := AppContext.Init()
	if err != nil {
		panic(err)
	}
}
