package main

import (
	"fmt"
	"reflect"
)

// WalkFunc is function that will be invoked during walking.
// Arguments:
// * value - currently walking value
// * level - nested level, starting with 0
// Return value - new value to update existing value, or NoUpdate{} if update not needed
type WalkFunc func(value interface{}, level int) (updateVal interface{})

// NoUpdate specifies that there's no update to existing value.
// It's distinguished from nil to allow setting nil to value.
type NoUpdate struct{}

// Walker is main struct for recursive walking.
// It contains walking options.
type Walker struct {
	// level is nested level, starting from 0
	level int

	// WalkFunc is function that will be invoked at every Walk invocation
	WalkFunc WalkFunc

	// NodeOnly will cause WalkFunc to be invoked only for node elements (struct, map, slice).
	// Self-exclusive with LeafOnly.
	NodeOnly bool

	// LeafOnly will cause WalkFunc to be invoked only for leaf elements.
	// Self-exclusive with NodeOnly.
	LeafOnly bool
}

func NewWalker(f WalkFunc) (*Walker, error) {
	if f == nil {
		return nil, fmt.Errorf("WalkFunc cannot be nil")
	}

	return &Walker{WalkFunc: f}, nil
}

// Go walks through nested object recursively and invokes WalkFunc.
// It looks inside structs, maps, slices.
func (w *Walker) Go(v interface{}) (ret interface{}) {
	ret = NoUpdate{}

	vv := reflect.ValueOf(v)
	kind := vv.Type().Kind()

	if kindOf(LeafKinds, kind) {
		if w.NodeOnly == false {
			ret = w.WalkFunc(v, w.level)
		}

		return ret
	}

	if kindOf(NodeKinds, kind) && w.LeafOnly == false {
		ret = w.WalkFunc(v, w.level)
		if !reflect.DeepEqual(ret, NoUpdate{}) {
			return ret
		}
	}

	w.level++
	defer func() { w.level-- }()

	switch kind {
	case reflect.Struct:
		num := vv.Type().NumField()
		for i := 0; i < num; i++ {
			if !vv.Field(i).CanInterface() {
				continue
			}

			ret := w.Go(vv.Field(i).Interface())
			if reflect.DeepEqual(ret, NoUpdate{}) {
				continue
			}

			if vv.Field(i).CanSet() {
				vv.Field(i).Set(reflect.ValueOf(ret))
			}
		}
	case reflect.Map:
		iter := vv.MapRange()
		for iter.Next() {
			//k := iter.Key()
			v := iter.Value()
			if !v.CanInterface() {
				continue
			}

			w.Go(v.Interface())
		}
	case reflect.Slice:
		n := vv.Len()
		for i := 0; i < n; i++ {
			elem := vv.Index(i)
			if !elem.CanInterface() {
				continue
			}

			w.Go(elem.Interface())
		}
	}

	return NoUpdate{}
}

var (
	LeafKinds = []reflect.Kind{reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128,
		reflect.Func,
		reflect.Interface,
		reflect.String}

	NodeKinds = []reflect.Kind{reflect.Map,
		reflect.Slice,
		reflect.Struct}

	PointerKind = []reflect.Kind{reflect.Ptr}
)

func kindOf(kinds []reflect.Kind, kind reflect.Kind) bool {
	for _, k := range kinds {
		if kind == k {
			return true
		}
	}
	return false
}
