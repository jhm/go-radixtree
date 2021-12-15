# go-radixtree [![Actions Status](https://github.com/jhm/go-radixtree/workflows/Main/badge.svg)](https://github.com/jhm/go-radixtree/actions)

An implementation of a mutable radix tree that uses byte slices for keys.
Insertion, deletion and searching operations all have a worst case of O(n) where
n is the length of the longest key in the tree. This implementation is not
thread safe.

## Basic Usage

```golang
// Creating a radix tree.
t := radixtree.New()

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

Upon the release of Go 1.18 the radix tree and its API will make use of generic
type parameters and the tests will most likely be updated to use the new fuzzing
API.
