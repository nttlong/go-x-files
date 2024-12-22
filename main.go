package main

import (
	"fmt"

	application "github.com/unvs/libs/application"
	"github.com/unvs/libs/cachery"
)

// decalre test strucutre
type testStruct struct {
	Name string
	Age  int
}

func main() {
	mytest := testStruct{Name: "John", Age: 30}
	//defer logger.HandlePanic()
	application.InitGlobalContext()
	cachery.Init(
		application.AppContext.Config.CacheServer,
		application.AppContext.Config.CachePrefix)
	cachery.HealthCheck()
	cachery.Set("testStruct", mytest, 0)
	cachery.Set("testint", 1222, 0)
	var testValue testStruct

	if cachery.Get("testStruct", &testValue) {
		fmt.Println(testValue)
	}
	var testint int
	if cachery.Get("testint", &testint) {
		fmt.Println(testint)
	}

	//logger.SetupLoggers(application.AppContext.AppPath, application.AppContext.Config.Log.Path)

	// app := application.AppContext

	// app.Cacher.SetText("test", "Hello, World!", cacher.WithExpiry(10))
	// test := app.Cacher.GetText("test")
	// app.Cacher.SetStruct("testStruct", mytest, cacher.WithExpiry(10))
	// mytest2 := testStruct{}
	// app.Cacher.GetStruct("testStruct", &mytest2)
	// //cast mytest2 to testStruct

	// fmt.Println(mytest2.Name, mytest2.Age)

	// fmt.Println(test)

}
