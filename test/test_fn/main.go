package main

import (
	"fmt"
	"reflect"
)

type Test struct {
	Name string
	Age  int
}

func test(t *Test) bool {
	return t.Name == "1234"
}

func GetBody(f interface{}) (string, error) {
	// Get the function value using reflection
	funcValue := reflect.ValueOf(f)
	if funcValue.Kind() != reflect.Func {
		return "", fmt.Errorf("input is not a function")
	}

	// Get function information
	funcType := funcValue.Type()
	numArgs := funcType.NumIn()
	numReturns := funcType.NumOut()

	// Build the function signature string
	signature := fmt.Sprintf("func(")
	for i := 0; i < numArgs; i++ {
		argType := funcType.In(i).Name()
		if i > 0 {
			signature += ", "
		}
		signature += argType
	}
	signature += ") ("

	// Add return types if any
	for i := 0; i < numReturns; i++ {
		returnType := funcType.Out(i).Name()
		if i > 0 {
			signature += ", "
		}
		signature += returnType
	}
	signature += ")"

	// We can't directly access the function body at runtime in Go
	body := "Function body not directly accessible at runtime"

	return signature + " {\n" + body + "\n}", nil
}

func main() {
	body, err := GetBody(test)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(body)
}
