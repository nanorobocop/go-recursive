package main

import (
	"reflect"
)

// WalkFunc is function that will be invoked during walking.
// Arguments:
// * value - currently walking value
// * level - nested level, starting with 0
// Return value - new value to update existing value, or NoUpdate{} if update not needed
type WalkFunc func(value interface{}, level int) (ret interface{})

// NoUpdate specifies that there's no update to existing value.
// It's distinguished from nil to allow setting nil as value.
type NoUpdate struct{}

var emptyValue = reflect.Value{}

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

// Go walks through nested object recursively and invokes WalkFunc.
// It looks inside structs, maps, slices.
// In order to update object, must call by pointer.
func Go(obj interface{}, f WalkFunc) {
	val := reflect.ValueOf(obj)
	w := &Walker{
		WalkFunc: f,
	}
	w.GoValue(val)
}

// GoValue walks through nested objects recursively.
// Returns updated value and changed flag.
func (w *Walker) GoValue(elem reflect.Value) (reflect.Value, bool) {
	kind := elem.Kind()

	elem2 := elem

	for {
		if kind == reflect.Interface {
			elem2 = elem2.Elem()
			kind = elem2.Kind()
			if elem2.CanAddr() {
				elem = elem2
			}
		} else {
			break
		}
	}

	if kindOf(LeafKinds, kind) && w.NodeOnly == false {
		orig := elem.Interface()
		ret := w.WalkFunc(orig, w.level)
		if !reflect.DeepEqual(orig, ret) {
			if elem.CanSet() {
				elem.Set(reflect.ValueOf(ret))
			} else {
				return reflect.ValueOf(ret), true
			}
		}
		return emptyValue, false
	}

	if kindOf(NodeKinds, kind) && w.LeafOnly == false {
		orig := elem.Interface()
		ret := w.WalkFunc(elem.Interface(), w.level)
		if !reflect.DeepEqual(orig, ret) {
			return emptyValue, true
		}
	}

	switch kind {
	case reflect.Struct:
		w.level++
		defer func() { w.level-- }()

		var newElemPtr reflect.Value

		num := elem.NumField()
		for i := 0; i < num; i++ {
			val := elem.Field(i)

			if !val.CanInterface() {
				continue
			}

			ret, changed := w.GoValue(val)
			if !changed {
				continue
			}

			if val.CanSet() {
				val.Set(ret)
				continue
			}

			if newElemPtr == emptyValue {
				newElemPtr = copyStruct(elem)
			}

			newElemPtr.Elem().Field(i).Set(ret)
		}
		if newElemPtr != emptyValue {
			return newElemPtr.Elem(), true
		}
		return emptyValue, false
	case reflect.Map:
		w.level++
		defer func() { w.level-- }()
		iter := elem.MapRange()
		for iter.Next() {
			key := iter.Key()
			val := iter.Value()

			if !val.CanInterface() {
				continue
			}

			ret, changed := w.GoValue(val)
			if !changed {
				continue
			}

			elem.SetMapIndex(key, ret)
		}
	case reflect.Slice:
		w.level++
		defer func() { w.level-- }()
		n := elem.Len()
		for i := 0; i < n; i++ {
			val := elem.Index(i)
			if !val.CanInterface() {
				continue
			}

			ret, changed := w.GoValue(val)
			if !changed {
				continue
			}

			val.Set(ret)
		}
	case reflect.Ptr:
		if elem.IsZero() {
			break
		}
		val := elem.Elem()
		w.GoValue(val)
	}

	return reflect.Value{}, false
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

// copyStruct copies struct to new one.
// Returns pointer to new struct.
// Used when original struct is not addressable.
func copyStruct(v reflect.Value) reflect.Value {
	res := reflect.New(v.Type())

	num := v.NumField()
	for i := 0; i < num; i++ {
		val := v.Field(i)
		res.Elem().Field(i).Set(val)
	}
	return res
}
