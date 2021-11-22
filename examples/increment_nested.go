package main

import (
	"fmt"
	"reflect"

	recursive "github.com/nanorobocop/go-recursive"
)

type NestedStruct struct {
	Int    int
	Nested *NestedStruct
}

func incrementFunc(obj interface{}, level int) interface{} {
	if i, ok := obj.(int); ok {
		i += +1
		return i
	}

	return obj // return unchanged
}

func printFunc(obj interface{}, level int) interface{} {
	val := reflect.ValueOf(obj)
	kind := val.Type().Kind()
	fmt.Printf("%s%+v (%s)\n", indent(level), val, kind)
	return obj // return unchanged
}

func indent(level int) (indent string) {
	if level == 0 {
		return ""
	}

	for i := 0; i < level; i++ {
		indent += ">>>>"
	}
	indent += " "
	return
}

func main() {
	nested := NestedStruct{Int: 20}
	object := map[int]NestedStruct{
		1: NestedStruct{Int: 10, Nested: &nested},
		2: NestedStruct{},
	}

	fmt.Println("Before:")
	recursive.Go(&object, printFunc)

	recursive.Go(&object, incrementFunc)

	fmt.Println()
	fmt.Println("After:")
	recursive.Go(&object, printFunc)
}
