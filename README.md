# Go-recursive - Walk and update nested Go objects

Go-recursive is Golang library to walk through and update nested objects.
It goes down structs, maps, slices recursively.

## Example

Incrememt Int value inside nested struct:

1. Define function:

   ```go
   func incrementFunc(obj interface{}, level int) interface{} {
   	if i, ok := obj.(int); ok {
   		i += +1
   		return i
   	}
   
   	return obj // return unchanged
   }
   ```

1. Define object. Let it be complex object here - nested structs inside map:

   ```go
   type NestedStruct struct {
   	Int    int
   	Nested *NestedStruct
   }
   
   var(
   	n = NestedStruct{Int: 20}
   	obj = map[int]NestedStruct{
   		1: NestedStruct{Int: 10, Nested: &n},
   		2: NestedStruct{},
   	}
   )
   ```

1. Execute:

   ```go
   	recursive.Go(&obj, incrementFunc)
   ```

1. Check result:

   ```txt
   Before:
   map[1:{Int:10 Nested:0xc000010230} 2:{Int:0 Nested:<nil>}] (map)
   >>>> {Int:0 Nested:<nil>} (struct)
   >>>>>>>> 0 (int)
   >>>>>>>> <nil> (ptr)
   >>>> {Int:10 Nested:0xc000010230} (struct)
   >>>>>>>> 10 (int)
   >>>>>>>> &{Int:20 Nested:<nil>} (ptr)
   >>>>>>>>>>>> {Int:20 Nested:<nil>} (struct)
   >>>>>>>>>>>>>>>> 20 (int)
   >>>>>>>>>>>>>>>> <nil> (ptr)

   After:
   map[1:{Int:11 Nested:0xc000010230} 2:{Int:1 Nested:<nil>}] (map)
   >>>> {Int:1 Nested:<nil>} (struct)
   >>>>>>>> 1 (int)
   >>>>>>>> <nil> (ptr)
   >>>> {Int:11 Nested:0xc000010230} (struct)
   >>>>>>>> 11 (int)
   >>>>>>>> &{Int:21 Nested:<nil>} (ptr)
   >>>>>>>>>>>> {Int:21 Nested:<nil>} (struct)
   >>>>>>>>>>>>>>>> 21 (int)
   >>>>>>>>>>>>>>>> <nil> (ptr)
   ```

Check [examples](examples) directory for complete program.
