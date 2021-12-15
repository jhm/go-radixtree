# go-radixtree

An implementation of a mutable radix tree that uses byte slices for keys.
Insertion, deletion and searching operations all have a worst case of O(n) where
n is the length of the longest key in the tree. This implementation is not
thread safe.

The main branch now requires Go 1.18 because the radix tree makes use of generic
type parameters. For a version that works on Go 1.17 and below see the v1.0.0
tag.

## Basic Usage

```golang
// Creating a radix tree.
t := radixtree.New[int]()

// Insert some values.
t.Insert([]byte{192, 1}, 0)
t.Insert([]byte{192, 1, 1, 4}, 1)

// Get the number of elements in the tree.
t.Len()

// Fetch a value.
v, found := t.Get([]byte{192, 1})

// Find all values that have a key that starts with 192.
vs := t.Find([]byte{192})

// Remove a value.
oldValue, found := t.Remove([]byte{192, 1})
```

See the Godocs for the rest of the API.

## Future Changes

The tests will most likely be updated to use the new fuzzing API released in Go
1.18.
