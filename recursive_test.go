package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

var (
	w   = &bytes.Buffer{}
	str = "str"
)

type SimpleStruct struct {
	B bool
	I int
	S string
	X *string
}

type SimpleStruct2 struct {
	I int
}

type ComplexStruct struct {
	SimpleStruct               // embedded struct
	A            SimpleStruct2 // nested struct
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

func printWalkFunc(obj interface{}, level int) (ret interface{}) {
	val := reflect.ValueOf(obj)
	kind := val.Type().Kind()
	fmt.Fprintf(w, "%s%+v (%s)\n", indent(level), val, kind)
	return obj
}

func TestPrintInt(t *testing.T) {
	value := 12
	expected := `12 (int)
`
	w = &bytes.Buffer{}
	Go(&value, printWalkFunc)
	a, _ := ioutil.ReadAll(w)

	if string(a) != expected {
		t.Errorf("Expected: %+v\nActual: %+v", expected, string(a))
	}
}

func TestPrintStruct(t *testing.T) {
	value := SimpleStruct{B: false, I: 5, S: "str", X: nil}
	expected := `{B:false I:5 S:str X:<nil>} (struct)
>>>> false (bool)
>>>> 5 (int)
>>>> str (string)
`
	Go(&value, printWalkFunc)
	a, _ := ioutil.ReadAll(w)

	if string(a) != expected {
		t.Errorf("Expected: %+v\nActual: %+v", expected, string(a))
	}
}

func TestPrintMap(t *testing.T) {
	value := map[string]int{"a": 11, "b": 12}
	expected := `map[a:11 b:12] (map)
>>>> 11 (int)
>>>> 12 (int)
`

	Go(&value, printWalkFunc)
	a, _ := ioutil.ReadAll(w)

	if string(a) != expected {
		t.Errorf("Expected: %+v\nActual: %+v", expected, string(a))
	}
}

func TestPrintSlice(t *testing.T) {
	value := []int{1, 2, 3, 4, 5}
	expected := `[1 2 3 4 5] (slice)
>>>> 1 (int)
>>>> 2 (int)
>>>> 3 (int)
>>>> 4 (int)
>>>> 5 (int)
`
	Go(&value, printWalkFunc)
	a, _ := ioutil.ReadAll(w)

	if string(a) != expected {
		t.Errorf("Expected: %+v\nActual: %+v", expected, string(a))
	}
}

func TestPrintComplex(t *testing.T) {
	value := ComplexStruct{
		SimpleStruct: SimpleStruct{
			B: true,
			I: 5,
			S: "str",
			X: &str,
		},
		A: SimpleStruct2{
			I: 10,
		},
	}
	expected := `{SimpleStruct:{B:true I:5 S:str X:0x11f7460} A:{I:10}} (struct)
>>>> {B:true I:5 S:str X:0x11f7460} (struct)
>>>>>>>> true (bool)
>>>>>>>> 5 (int)
>>>>>>>> str (string)
>>>>>>>> str (string)
>>>> {I:10} (struct)
>>>>>>>> 10 (int)
`
	Go(&value, printWalkFunc)
	a, _ := ioutil.ReadAll(w)

	if string(a) != expected {
		t.Errorf("Expected: %+v\nActual: %+v", expected, string(a))
	}
}

func incIntWalkFunc(v interface{}, l int) interface{} {
	if i, ok := v.(int); ok {
		i += +1
		return i
	}

	return v
}

func TestIncInt(t *testing.T) {
	value := 5
	expected := 6

	Go(&value, incIntWalkFunc)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected: %+v, actual: %+v", expected, value)
	}
}

func TestIncIntWithString(t *testing.T) {
	value := "asdf"
	expected := "asdf"

	Go(&value, incIntWalkFunc)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected: %+v, actual: %+v", expected, value)
	}
}

func TestIncIntWithStruct(t *testing.T) {
	value := SimpleStruct{B: false, I: 5, S: "str"}
	expected := SimpleStruct{B: false, I: 6, S: "str"}

	Go(&value, incIntWalkFunc)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected: %+v, actual: %+v", expected, value)
	}
}

func TestIncIntWithMap(t *testing.T) {
	value := map[string]int{"a": 5, "b": 10}
	expected := map[string]int{"a": 6, "b": 11}

	Go(&value, incIntWalkFunc)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected: %+v, actual: %+v", expected, value)
	}
}

func TestIncIntWithSlice(t *testing.T) {
	value := []int{5, 10}
	expected := []int{6, 11}

	Go(&value, incIntWalkFunc)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected: %+v, actual: %+v", expected, value)
	}
}

func TestIncIntWithComplexStruct(t *testing.T) {
	value := ComplexStruct{
		SimpleStruct: SimpleStruct{I: 10},
		A:            SimpleStruct2{I: 12},
	}

	expected := ComplexStruct{
		SimpleStruct: SimpleStruct{I: 11},
		A:            SimpleStruct2{I: 13},
	}

	Go(&value, incIntWalkFunc)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected: %+v, actual: %+v", expected, value)
	}
}

func TestIncIntWithMapComplexStruct(t *testing.T) {
	value := map[int]ComplexStruct{
		1: ComplexStruct{
			SimpleStruct: SimpleStruct{I: 10},
			A:            SimpleStruct2{I: 12},
		},
		2: ComplexStruct{
			SimpleStruct: SimpleStruct{I: 20},
			A:            SimpleStruct2{I: 22},
		},
	}
	expected := map[int]ComplexStruct{
		1: ComplexStruct{
			SimpleStruct: SimpleStruct{I: 11},
			A:            SimpleStruct2{I: 13},
		},
		2: ComplexStruct{
			SimpleStruct: SimpleStruct{I: 21},
			A:            SimpleStruct2{I: 23},
		},
	}
	Go(&value, incIntWalkFunc)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected: %+v, actual: %+v", expected, value)
	}
}

func updateStruct(i interface{}, level int) interface{} {
	s, ok := i.(SimpleStruct)
	if !ok {
		return i
	}

	s.I += 1
	return s
}

func TestUpdateStruct(t *testing.T) {
	s := map[int]SimpleStruct{
		1: SimpleStruct{B: true, I: 1, S: "s", X: &str},
		2: SimpleStruct{},
	}

	expected := map[int]SimpleStruct{
		1: SimpleStruct{B: true, I: 2, S: "s", X: &str},
		2: SimpleStruct{I: 1},
	}

	Go(&s, updateStruct)

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected: %+v\nActual: %+v", expected, s)
	}
}

func TestCopyStruct(t *testing.T) {
	s := SimpleStruct{B: true, I: 1, S: "s", X: &str}

	s2 := CopyStruct(reflect.ValueOf(s)).Interface().(*SimpleStruct)

	if !reflect.DeepEqual(s, *s2) {
		t.Errorf("Expected: %+v\nActual: %+v", s, *s2)
	}
}
