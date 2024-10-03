
# Intro to Go Lab 1: Imperative Programming - Solutions

### Question 3b

Arrays are always passed by value. Slices, although passes by value, are effectively pointers to arrays. More specifically, a slice is a value storing a pointer to an underlying array. It also stores length and capacity of the array.

To solve the problem change the signature of `mapArray` to:

```go
func mapArray(f func(a int) int, array [3]int) [3]int
```

Alternatively, you could pass pointers to arrays [but this is not recommended](https://golang.org/doc/effective_go.html#arrays).

### Question 3c

`mapArray` no longer works because the type `[5]int` doesn't match type `[3]int`. This shows a fundamental issue with arrays and makes them much more limited than they are in C.

### Question 3d

```go
newSlice := intsSlice[1:3] // = [3, 4]
mapSlice(square, newSlice)
fmt.Println(newSlice)  // = [9, 16]
fmt.Println(intsSlice) // = [2 9 16 5 6]
```

This happens because the slice returned by the slice operator (`:`) returns a new slice which still holds a pointer to the same underlying array.

### Question 3e

Append may reallocate the slice and return the newly allocated one. Therefore you have to return it from `double`:

```go
func double(slice []int) []int {
    slice = append(slice, slice...)
    return slice
}
```

A slice is basically the following struct:

```go
type slice struct {
    length int
    capacity int
    pointerToArray *[capacity]elementType
}
```

When we pass a slice to a function, we pass a copy of this "struct" (by value).

Append works roughly in the following way: (in pseudocode, in reality far more complicated)

```go
neededCapacity := slice.length + length(elementsToAppend)
if slice.capacity < neededCapacity {
    slice.capacity = slice.capacity * 2
    var newArray[slice.capacity]
    copy(newArray, oldArray)
    slice.pointerToArray = &newArray
}
actuallyAppend()
...
```

In words:

After appending, the slice that you have in the `double()` function will be storing a pointer to a completely new array with a copy of the previous values and then the appended ones. The append operation is not performed on the old array and that's why the array from the slice in `main()` does not get updated.

In summary, the only pointer we are dealing with is the pointer to the array inside the struct. The struct itself is passed by value. Thus, any operation that modifies the length or capacity requires you to reassign the slice.

### Question 3f

- Both are passed by value. When passing an array to a function you pass a copy of the actual values. When passing a slice you pass a run-time data structure (often 24 bytes). It contains a pointer to an array, its length and its capacity. Hence slices can be seen as an abstraction of arrays.
- Append appends the elements to the end of the slice and returns the result. The result needs to be returned because the underlying array may change if the length of the final array would be greater than the current array's capacity.
- Arrays have very limited use cases. 99% of the time you will want to use a slice.
