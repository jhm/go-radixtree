// Package radixtree provides an implementation of a mutable radix tree.
// Insertion, deletion and searching operations all have a worst case of O(n)
// where n is the length of the longest key in the tree. This implementation is
// not thread safe.
package radixtree

import (
	"bytes"
	"sort"
)

// children encapsulates a slice of nodes sorted in ascending order by the first
// byte of their prefix.
type children[T any] []*node[T]

func (c *children[T]) add(node *node[T]) {
	i := c.search(node.prefix[0])
	*c = append(*c, nil)
	copy((*c)[i+1:], (*c)[i:])
	(*c)[i] = node
}

func (c children[T]) get(b byte) *node[T] {
	if i := c.index(b); i >= 0 {
		return c[i]
	}
	return nil
}

func (c children[T]) index(b byte) int {
	if i := c.search(b); i < len(c) && c[i].prefix[0] == b {
		return i
	}
	return -1
}

func (c children[T]) search(b byte) int {
	return sort.Search(len(c), func(i int) bool {
		return c[i].prefix[0] >= b
	})
}

// node encapsulates a prefix, with a possible associated value, and a set of
// child nodes.
type node[T any] struct {
	prefix   []byte
	children children[T]
	value    *T
}

func (n *node[T]) hasValue() bool {
	return n.value != nil
}

func (n *node[T]) max() (T, bool) {
	for len(n.children) > 0 {
		n = n.children[len(n.children)-1]
	}
	if n.hasValue() {
		return *n.value, true
	}
	var zero T
	return zero, false
}

func (n *node[T]) min() (T, bool) {
	for !n.hasValue() && len(n.children) > 0 {
		n = n.children[0]
	}
	if n.hasValue() {
		return *n.value, true
	}
	var zero T
	return zero, false
}

// RadixTree implements a mutable radix tree.
type RadixTree[T any] struct {
	root *node[T]
	size int
}

// New creates and returns an empty radix tree.
func New[T any]() *RadixTree[T] {
	return &RadixTree[T]{root: &node[T]{}}
}

// Contains returns true if key is in the tree, false otherwise.
func (t *RadixTree[T]) Contains(key []byte) bool {
	_, b := t.Get(key)
	return b
}

// Find returns a slice that contains all of the values that have a key that
// starts with the given prefix. The slice will be ordered in ascending key
// order.
func (t *RadixTree[T]) Find(prefix []byte) []T {
	var results []T
	t.Walk(prefix, func(value T) bool {
		results = append(results, value)
		return true
	})
	return results
}

// Get returns the value associated with the given key. If the key is found in
// the tree it returns the associated value and a boolean value of true
// indicating that a value was found. If the key is not in the tree it returns
// the zero value for type T and a false boolean value.
func (t *RadixTree[T]) Get(key []byte) (T, bool) {
	n := t.root

	for len(key) > 0 {
		n = n.children.get(key[0])
		if n == nil || !bytes.HasPrefix(key, n.prefix) {
			var zero T
			return zero, false
		}
		key = key[len(n.prefix):]
	}

	if n.hasValue() {
		return *n.value, true
	}
	var zero T
	return zero, false
}

// Insert adds the value to the radix tree with the given key. If the exact key
// already exists in the radix tree it updates the value and returns the old
// value and a boolean value of true indicating that an old value was found. If
// the key was not in the tree it returns the zero value for type T and a false
// boolean value.
func (t *RadixTree[T]) Insert(key []byte, value T) (T, bool) {
	n := t.root

	for len(key) > 0 {
		i := n.children.index(key[0])
		if i < 0 {
			// There is no child starting with the first byte of the
			// key so we can simply add a new child node to n.
			n.children.add(&node[T]{value: &value, prefix: key})
			t.size++
			var zero T
			return zero, false
		}

		child := n.children[i]
		lcm := longestCommonPrefix(key, child.prefix)
		if lcm < len(child.prefix) {
			// The child needs to be split.
			newChild := &node[T]{prefix: key[:lcm]}
			n.children[i] = newChild
			child.prefix = child.prefix[lcm:]
			newChild.children.add(child)
			key = key[lcm:]
			if len(key) == 0 {
				newChild.value = &value
				t.size++
				var zero T
				return zero, false
			}
			newChild.children.add(&node[T]{value: &value, prefix: key})
			t.size++
			var zero T
			return zero, false
		}
		n = child
		key = key[lcm:]
	}

	if n.hasValue() {
		// This insert is actually an update to an existing value.
		old := *n.value
		n.value = &value
		return old, true
	}
	// The node exists but doesn't contain a value.
	n.value = &value
	t.size++
	var zero T
	return zero, false
}

// Len returns the number of values in the tree.
func (t *RadixTree[T]) Len() int {
	return t.size
}

// LongestPrefix returns the value associated with the key that has the longest
// prefix of the given key. If a value is found it returns the value and a
// boolean value of true. If no value is found it returns the zero value for
// type T and a boolean value of false.
func (t *RadixTree[T]) LongestPrefix(key []byte) (T, bool) {
	n := t.root
	var last *T

	for len(key) > 0 {
		n = n.children.get(key[0])
		if n == nil || !bytes.HasPrefix(key, n.prefix) {
			break
		}
		if n.hasValue() {
			last = n.value
		}
		key = key[len(n.prefix):]
	}
	if last != nil {
		return *last, true
	}
	var zero T
	return zero, false
}

// Max returns the value associated with the largest key in the tree. The
// boolean return value will be true if a maximum value was found and false if
// the tree is empty and therefore has no maximum value.
func (t *RadixTree[T]) Max() (T, bool) {
	return t.root.max()
}

// Min returns the value associated with the smallest key in the tree. The
// boolean return value will be true if a maximum value was found and false if
// the tree is empty and therefore has no minimum value.
func (t *RadixTree[T]) Min() (T, bool) {
	return t.root.min()
}

// Predecessor returns the value that is associated with the key that
// immediately precedes the given key. If a predecessor is found, its value and
// a boolean value of true will returned. If there is no predecessor, or the
// given key does not exist in the tree, the zero value for type T and a boolean
// value of false will be returned.
func (t *RadixTree[T]) Predecessor(key []byte) (T, bool) {
	n := t.root
	ancestor := false
	var min *node[T]

	for len(key) > 0 {
		i := n.children.index(key[0])
		if i < 0 || !bytes.HasPrefix(key, n.children[i].prefix) {
			var zero T
			return zero, false
		}

		if i > 0 {
			min = n.children[i-1]
			ancestor = false
		} else if n.hasValue() {
			min = n
			ancestor = true
		}
		n = n.children[i]
		key = key[len(n.prefix):]
	}

	if min != nil {
		if ancestor {
			return *min.value, true
		}
		return min.max()
	}
	var zero T
	return zero, false
}

// Remove removes the key and its associated value from the tree and returns the
// old value and a boolean value of true indicating that the given key was
// found. If the key was not present in the tree it will return the zero value
// for type T and a boolean value of false.
func (t *RadixTree[T]) Remove(key []byte) (T, bool) {
	var parent *node[T]
	var i int
	n := t.root
	root := n

	for len(key) > 0 {
		if i = n.children.index(key[0]); i < 0 {
			var zero T
			return zero, false
		}
		parent = n
		n = n.children[i]
		if !bytes.HasPrefix(key, n.prefix) {
			var zero T
			return zero, false
		}
		key = key[len(n.prefix):]
	}

	if n.hasValue() {
		v := *n.value
		n.value = nil

		// If the node to be deleted has no children it can be removed
		// from the parent node's list of children.
		if parent != nil && len(n.children) == 0 {
			parent.children = append(parent.children[:i], parent.children[i+1:]...)
		}

		// If the node to be deleted only has a single child that child
		// can be merged into node n.
		if n != root && len(n.children) == 1 {
			merge(n)
		}

		// If the parent node exists, has no value, and only has a
		// single child it can be merged with that child.
		if parent != nil && parent != root && len(parent.children) == 1 && !parent.hasValue() {
			merge(parent)
		}
		t.size--
		return v, true
	}
	var zero T
	return zero, false
}

func merge[T any](n *node[T]) {
	child := n.children[0]
	n.prefix = append(n.prefix, child.prefix...)
	n.value = child.value
	n.children = child.children
}

// Successor returns the value that is associated with the key that immediately
// follows the given key. If a successor is found, its value and a boolean value
// of true will be returned. If there is no successor, or the given key does not
// exist in the tree, the zero value for type T and a boolean value of false
// will be returned.
func (t *RadixTree[T]) Successor(key []byte) (T, bool) {
	n := t.root
	var min *node[T]

	for len(key) > 0 {
		i := n.children.index(key[0])
		if i < 0 || !bytes.HasPrefix(key, n.children[i].prefix) {
			var zero T
			return zero, false
		}
		if r := i + 1; r < len(n.children) {
			min = n.children[r]
		}
		n = n.children[i]
		key = key[len(n.prefix):]
	}

	if len(n.children) != 0 {
		min = n.children[0]
	}

	if min != nil {
		return min.min()
	}
	var zero T
	return zero, false
}

// Values returns all of the values in the tree in the ascending order of their
// keys.
func (t *RadixTree[T]) Values() []T {
	results := make([]T, 0, t.Len())
	t.Walk([]byte{}, func(value T) bool {
		results = append(results, value)
		return true
	})
	return results
}

// Walk traverses the tree rooted at the given prefix and executes function f
// for each value. If f returns true the traversal continues otherwise the
// traversal stops.
func (t *RadixTree[T]) Walk(prefix []byte, f func(value T) bool) {
	n := t.root

	for len(prefix) > 0 {
		n = n.children.get(prefix[0])
		if n == nil || !bytes.HasPrefix(prefix, n.prefix) {
			break
		}
		prefix = prefix[len(n.prefix):]
	}

	if n != nil {
		walk(n, f)
	}
}

func walk[T any](n *node[T], f func(value T) bool) bool {
	if n.hasValue() && !f(*n.value) {
		return false
	}
	for _, node := range n.children {
		if !walk(node, f) {
			return false
		}
	}
	return true
}

func longestCommonPrefix(a, b []byte) int {
	limit := len(a)
	if l := len(b); l < limit {
		limit = l
	}

	i := 0
	for ; i < limit; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return i
}
