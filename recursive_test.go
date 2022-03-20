package recursive

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"testing"
)

var (
	w   = &bytes.Buffer{}
	str = "str"
)

type SimpleStruct struct {
	Bool bool
	Int  int
	Str  string
	Ptr  *string
}

type SimpleStruct2 struct {
	Int int
}

type ComplexStruct struct {
	SimpleStruct               // embedded struct
	A            SimpleStruct2 // nested struct
}

type NestedStruct struct {
	Int    int
	Nested *NestedStruct
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
	return obj // return unchanged
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
	value := SimpleStruct{Bool: false, Int: 5, Str: "str", Ptr: nil}
	expected := `{Bool:false Int:5 Str:str Ptr:<nil>} (struct)
>>>> false (bool)
>>>> 5 (int)
>>>> str (string)
>>>> <nil> (ptr)
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
			Bool: true,
			Int:  5,
			Str:  "str",
			Ptr:  &str,
		},
		A: SimpleStruct2{
			Int: 10,
		},
	}
	expected := `{SimpleStruct:{Bool:true Int:5 Str:str Ptr:0x[0-9a-f]+} A:{Int:10}} \(struct\)
>>>> {Bool:true Int:5 Str:str Ptr:0x[0-9a-f]+} \(struct\)
>>>>>>>> true \(bool\)
>>>>>>>> 5 \(int\)
>>>>>>>> str \(string\)
>>>>>>>> 0x[0-9a-f]+ \(ptr\)
>>>>>>>>>>>> str \(string\)
>>>> {Int:10} \(struct\)
>>>>>>>> 10 \(int\)
`
	Go(&value, printWalkFunc)
	a, _ := ioutil.ReadAll(w)

	match, _ := regexp.MatchString(expected, string(a))
	if !match {
		t.Errorf("Expected: %+v\nActual: %+v", expected, string(a))
	}
}

func incIntWalkFunc(v interface{}, l int) interface{} {
	if i, ok := v.(int); ok {
		i += +1
		return i
	}

	return v // return unchanged
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
	value := SimpleStruct{Bool: false, Int: 5, Str: "str"}
	expected := SimpleStruct{Bool: false, Int: 6, Str: "str"}

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
		SimpleStruct: SimpleStruct{Int: 10},
		A:            SimpleStruct2{Int: 12},
	}

	expected := ComplexStruct{
		SimpleStruct: SimpleStruct{Int: 11},
		A:            SimpleStruct2{Int: 13},
	}

	Go(&value, incIntWalkFunc)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected: %+v, actual: %+v", expected, value)
	}
}

func TestIncIntWithMapComplexStruct(t *testing.T) {
	value := map[int]ComplexStruct{
		1: ComplexStruct{
			SimpleStruct: SimpleStruct{Int: 10},
			A:            SimpleStruct2{Int: 12},
		},
		2: ComplexStruct{
			SimpleStruct: SimpleStruct{Int: 20},
			A:            SimpleStruct2{Int: 22},
		},
	}
	expected := map[int]ComplexStruct{
		1: ComplexStruct{
			SimpleStruct: SimpleStruct{Int: 11},
			A:            SimpleStruct2{Int: 13},
		},
		2: ComplexStruct{
			SimpleStruct: SimpleStruct{Int: 21},
			A:            SimpleStruct2{Int: 23},
		},
	}
	Go(&value, incIntWalkFunc)
	if !reflect.DeepEqual(value, expected) {
		t.Errorf("Expected: %+v, actual: %+v", expected, value)
	}
}

func TestNestedStruct(t *testing.T) {
	n := NestedStruct{Int: 20}
	s := map[int]NestedStruct{
		1: NestedStruct{Int: 10, Nested: &n},
		2: NestedStruct{},
	}

	expected := map[int]NestedStruct{
		1: NestedStruct{Int: 11, Nested: &NestedStruct{Int: 21}},
		2: NestedStruct{Int: 1},
	}

	Go(&s, incIntWalkFunc)

	if !reflect.DeepEqual(s, expected) || !reflect.DeepEqual(s[1].Nested, expected[1].Nested) {
		t.Errorf("Expected: %+v\nActual: %+v", expected, s)
		t.Errorf("Expected: %+v\nActual: %+v", expected[1].Nested, s[1].Nested)
	}
}

func updateStruct(i interface{}, level int) interface{} {
	s, ok := i.(SimpleStruct)
	if !ok {
		return i
	}

	s.Int += 1
	return s
}

func TestUpdateStruct(t *testing.T) {
	s := map[int]SimpleStruct{
		1: SimpleStruct{Bool: true, Int: 1, Str: "s", Ptr: &str},
		2: SimpleStruct{},
	}

	expected := map[int]SimpleStruct{
		1: SimpleStruct{Bool: true, Int: 2, Str: "s", Ptr: &str},
		2: SimpleStruct{Int: 1},
	}

	Go(&s, updateStruct)

	if !reflect.DeepEqual(s, expected) {
		t.Errorf("Expected: %+v\nActual: %+v", expected, s)
	}
}

func TestCopyStruct(t *testing.T) {
	s := SimpleStruct{Bool: true, Int: 1, Str: "s", Ptr: &str}

	s2 := CopyStruct(reflect.ValueOf(s)).Interface().(*SimpleStruct)

	if !reflect.DeepEqual(s, *s2) {
		t.Errorf("Expected: %+v\nActual: %+v", s, *s2)
	}
}
