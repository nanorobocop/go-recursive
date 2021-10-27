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

var w io.ReadWriter

type SimpleStruct struct {
	b bool
	i int
	s string
	x *string
}

type SimpleStruct2 struct {
	i int
}

type MultitypeStruct struct {
	SimpleStruct               // embedded struct
	a            SimpleStruct2 // nested struct
	b            bool
	i            int
	s            string
	x            *string
}

func PrintWalkFunc(value interface{}, level int) (updateVal interface{}) {
	vv := reflect.ValueOf(value)
	kind := vv.Type().Kind().String()
	out := fmt.Sprintf("%s%+v (%s)\n", indent(level), value, kind)
	fmt.Fprintf(w, out)
	return NoUpdate{}
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

func compare(t *testing.T, value interface{}, expected string) {
	w = &bytes.Buffer{}
	walker, _ := NewWalker(PrintWalkFunc)
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

func TestRecursiveInt(t *testing.T) {
	value := 12
	expected := `12 (int)
`

	compare(t, value, expected)
}

func TestRecursiveMap(t *testing.T) {
	value := map[string]int{"a": 11, "b": 12}
	expected := `map[a:11 b:12] (map)
>>>> 11 (int)
>>>> 12 (int)
`

	compare(t, value, expected)
}

func TestRecursiveSlice(t *testing.T) {
	value := []int{1, 2, 3, 4, 5}
	expected := `[1 2 3 4 5] (slice)
>>>> 1 (int)
>>>> 2 (int)
>>>> 3 (int)
>>>> 4 (int)
>>>> 5 (int)
`

	compare(t, value, expected)
}
