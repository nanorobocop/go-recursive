package main

import (
	"fmt"
	"reflect"
)

// WalkFunc is function that will be invoked during walking for each element.
// Arguments:
// level - nested level, starting 0
// value - currently walking value
// retValue - new value to update existing value
type WalkFunc func(level int, value interface{}) (retValue interface{})

// NoUpdate specifies that there's no update to existing value.
// It's distinguish from nil to allow setting nil to value.
type NoUpdate struct{}

func PrintWalkFunc(level int, value interface{}) (retValue interface{}) {
	vv := reflect.ValueOf(value)
	kind := vv.Type().Kind().String()
	fmt.Printf("%s%s\n", indent(level), kind)
	return NoUpdate{}
}

func PrintLeafFunc(level int, value interface{}) (retValue interface{}) {
	fmt.Printf("%s%v\n", indent(level), value)
	return NoUpdate{}
}

func indent(level int) (indent string) {
	for i := 0; i < level; i++ {
		indent += ">>>>"
	}
	return
}

// WalkOpts defines walking options.
// It could be empty struct at the beginning.
type WalkOpts struct {
	// Level is nested level, keep 0 at the beginning
	Level int

	// WalkFunc is function that will be invoked at every Walk invokaction
	WalkFunc WalkFunc

	// LeafFunc is func that will be invoked on for every leaf element (if not struct, map or slice)
	LeafFunc WalkFunc

	// NodeFunc is func that will be invoked only for node element (not struct, map or slice)
	NodeFunc WalkFunc
}

func (o *WalkOpts) Inc() *WalkOpts {
	o2 := &WalkOpts{}
	*o2 = *o
	o2.Level++
	return o2
}

// Walk walks through nested object recursively and invokes f.
// It looks inside structs, maps, slices.
func Walk(v interface{}, opts *WalkOpts) interface{} {
	vv := reflect.ValueOf(v)

	if opts.WalkFunc != nil {
		opts.WalkFunc(opts.Level, v)
	}

	switch vv.Type().Kind().String() {
	case "struct":
		num := vv.Type().NumField()
		for i := 0; i < num; i++ {
			if !vv.Field(i).CanInterface() {
				continue
			}
			Walk(vv.Field(i).Interface(), opts.Inc())
		}
	case "map":
		kind := vv.Type().Elem().Kind().String()
		if kind != "struct" &&
			kind != "map" &&
			kind != "slice" {
			return nil
		}

		iter := vv.MapRange()
		for iter.Next() {
			//k := iter.Key()
			v := iter.Value()
			Walk(v, opts.Inc())
		}
	case "slice":
		kind := vv.Type().Elem().Kind().String()
		if kind != "struct" &&
			kind != "map" &&
			kind != "slice" {
			return nil
		}

		n := vv.Len()
		for i := 0; i < n; i++ {
			Walk(vv.Index(i), opts.Inc())
		}
	default:
		opts.WalkFunc(opts.Level, v)
	}

	return nil
}
