package main

import (
	"fmt"

	ctx "github.com/unvs/libs/db/ctx"
	expr "github.com/unvs/libs/db/expr"
)

func main() {
	// fx :=expr.AnalyzeExpressionWithPlaceholders("Code==?","12345")
	fx, err := expr.GetMongoQueryFromString("endswith(Emp.Code,?)", "123")
	if err != nil {
		panic(err)
	}
	fmt.Println(expr.ToPrettyJSONOfBSON(fx))
	ctx_test, err := ctx.NewDBContext("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	db := ctx_test.GetDB("testdb")
	db.
	// use ast get body of test function

}
