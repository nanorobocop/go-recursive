package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
	//"os"
)

var (
	w   io.ReadWriter
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

type MultitypeStruct struct {
	SimpleStruct               // embedded struct
	A            SimpleStruct2 // nested struct
	B            bool
	I            int
	S            string
	X            *string
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

func printWalkFunc(value interface{}, level int) (updateVal interface{}) {
	vv := reflect.ValueOf(value)
	kind := vv.Type().Kind().String()
	out := fmt.Sprintf("%s%+v (%s)\n", indent(level), value, kind)
	fmt.Fprintf(w, out)
	return NoUpdate{}
}

func comparePrint(t *testing.T, value interface{}, expected string) {
	w = &bytes.Buffer{}
	walker, _ := NewWalker(printWalkFunc)
	walker.Go(value)
	a, _ := ioutil.ReadAll(w)
	actual := string(a)

	fail := false
	if actual != expected {
		fail = true
	}

	if testing.Verbose() || fail {
		t.Log("Value    :", value)
		t.Log("Expected :", expected)
		t.Log("Actual   :", actual)
	}

	if fail {
		t.Fail()
	}
}

func TestPrintInt(t *testing.T) {
	value := 12
	expected := `12 (int)
`

	comparePrint(t, value, expected)
}

func TestPrintStruct(t *testing.T) {
	value := SimpleStruct{B: false, I: 5, S: "str", X: nil}
	expected := `{B:false I:5 S:str X:<nil>} (struct)
>>>> false (bool)
>>>> 5 (int)
>>>> str (string)
`

	comparePrint(t, value, expected)
}

func TestPrintMap(t *testing.T) {
	value := map[string]int{"a": 11, "b": 12}
	expected := `map[a:11 b:12] (map)
>>>> 11 (int)
>>>> 12 (int)
`

	comparePrint(t, value, expected)
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

	comparePrint(t, value, expected)
}

func incIntWalkFunc(v interface{}, l int) interface{} {
	vv := reflect.ValueOf(v)
	kind := vv.Type().Kind()

	if kind == reflect.Int {
		i := v.(int)
		return i + 1
	}

	return NoUpdate{}
}

func TestIncInt(t *testing.T) {
	tests := []struct {
		title    string
		value    interface{}
		expected interface{}
	}{
		{
			title:    "simple int",
			value:    5,
			expected: 6,
		},
		{
			title:    "simple string",
			value:    "asdf",
			expected: NoUpdate{},
		},
		{
			title:    "simple bool",
			value:    false,
			expected: NoUpdate{},
		},
		{
			title:    "struct",
			value:    SimpleStruct{B: false, I: 5, S: "str"},
			expected: SimpleStruct{B: false, I: 6, S: "str"},
		},
	}

	walker, _ := NewWalker(incIntWalkFunc)

	for _, test := range tests {
		actual := walker.Go(test.value)

		if !reflect.DeepEqual(test.expected, actual) {
			t.Logf("Value    : %+v", test.value)
			t.Logf("Expected : %+v", test.expected)
			t.Logf("Actaul   : %+v", actual)
			t.Fail()
		}
	}

}
