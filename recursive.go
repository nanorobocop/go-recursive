package recursive

import (
	"reflect"
)

// WalkFunc is function that will be invoked during walking.
// Arguments:
// * value - currently walking value
// * level - nested level, starting with 0
// Return value - new value to update existing value, return original value if update not needed
type WalkFunc func(value interface{}, level int) (ret interface{})

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
			return reflect.ValueOf(ret), true
		}
	}

	switch kind {
	case reflect.Interface:
		elem = elem.Elem()
		ret, changed := w.GoValue(elem)
		if !changed {
			return emptyValue, false
		}

		if elem.CanSet() {
			elem.Set(ret)
			return emptyValue, false
		}
		return ret, true
	case reflect.Ptr:
		if elem.IsZero() {
			break
		}
		val := elem.Elem()
		w.GoValue(val)
	}

	w.level++
	defer func() { w.level-- }()

	switch kind {
	case reflect.Struct:
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
				newElemPtr = CopyStruct(elem)
			}

			newElemPtr.Elem().Field(i).Set(ret)
		}
		if newElemPtr != emptyValue {
			return newElemPtr.Elem(), true
		}
		return emptyValue, false
	case reflect.Map:
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
	}

	return emptyValue, false
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

// CopyStruct copies struct to new one.
// Returns pointer to new struct.
// Used when original struct is not addressable.
func CopyStruct(v reflect.Value) reflect.Value {
	res := reflect.New(v.Type())

	num := v.NumField()
	for i := 0; i < num; i++ {
		val := v.Field(i)
		res.Elem().Field(i).Set(val)
	}
	return res
}
