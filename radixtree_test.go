package radixtree

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func init() {
	sort.Strings(words)
}

func build(keys []string) *RadixTree[string] {
	tree := New[string]()
	for _, key := range keys {
		tree.Insert([]byte(key), key)
	}
	return tree
}

func hasPrefix(prefix string, xs []string) []string {
	var ys []string
	for _, s := range xs {
		if strings.HasPrefix(s, prefix) {
			ys = append(ys, s)
		}
	}
	return ys
}

func TestContains(t *testing.T) {
	tree := build(words)

	for _, key := range words {
		if !tree.Contains([]byte(key)) {
			t.Errorf("Contains(%s) returned false for an existing key", key)
		}
	}

	if tree.Contains([]byte{0}) {
		t.Errorf("Contains returned true for a non-existent key")
	}

	if tree.Contains([]byte{}) {
		t.Errorf("Contains returned true for empty byte slice key")
	}
}

func TestFind(t *testing.T) {
	tree := build(words)

	prefix := "t"
	want := hasPrefix(prefix, words)
	if got := tree.Find([]byte(prefix)); !reflect.DeepEqual(got, want) {
		t.Errorf("Find(%s)\n got: %v\nwant: %v", prefix, got, want)
	}

	prefix = "to"
	want = hasPrefix(prefix, words)
	if got := tree.Find([]byte(prefix)); !reflect.DeepEqual(got, want) {
		t.Errorf("Find(%s)\n got: %v\nwant: %v", prefix, got, want)
	}

	want = []string{}
	if got := tree.Find([]byte{0}); len(got) != 0 {
		t.Errorf("Find with a non-existent prefix\n got: %v\nwant: %v", got, want)
	}
}

func TestGet(t *testing.T) {
	tree := build(words)

	for _, want := range words {
		if got, ok := tree.Get([]byte(want)); got != want || !ok {
			t.Errorf("Get(%s)\n got: (%s, %t)\nwant: (%s, true)", want, got, ok, want)
		}
	}

	if got, ok := tree.Get([]byte{0}); ok || got != "" {
		t.Errorf("Get with a non-existent key\n got: (%v, %t)\nwant: (\"\", false)", got, ok)
	}

	// Insert a value with an existing key and get the updated value.
	want := "aardvark"
	tree.Insert([]byte(want), want)
	if got, ok := tree.Get([]byte(want)); !ok || got != want {
		t.Errorf("Get(%s) after insert with an existing key\n got: (%s, %t)\nwant: (%s, true)", want, got, ok, want)
	}

	// Delete a value that shares a prefix with another node and make sure
	// the other node still has the correct value.
	tree.Remove([]byte("to"))
	want = "toa"
	if got, ok := tree.Get([]byte(want)); !ok || got != want {
		t.Errorf("Get(%s) after deletion of sibling node\n got: (%s, %t)\nwant: (%s, true)", want, got, ok, want)
	}

	// Insert into an existing node that doesn't have a value.
	want = "aard"
	tree.Insert([]byte(want), want)
	if got, ok := tree.Get([]byte(want)); !ok || got != want {
		t.Errorf("Get(%s) after insert into a node without a value\n got: (%s, %t)\nwant: (%s, true)", want, got, ok, want)
	}
}

func TestInsert(t *testing.T) {
	tree := build(words)
	want := "wink"
	if got, ok := tree.Insert([]byte(want), want); got != want || !ok {
		t.Errorf("Insert to existing key\n got: (%s, %t)\nwant: (%s, true)", got, ok, want)
	}

	want = "wilt"
	tree.Insert([]byte(want), want)
	if !tree.Contains([]byte(want)) {
		t.Errorf("Contains(%s) false after split insert", want)
	}
}

func TestLen(t *testing.T) {
	tree := New[int]()
	if got := tree.Len(); got != 0 {
		t.Errorf("Len on empty tree\n got: %d\nwant: 0", got)
	}

	for i, k := range words {
		tree.Insert([]byte(k), i)
		want := i + 1
		if got := tree.Len(); got != want {
			t.Errorf("Len\n got: %d\nwant: %d", got, want)
		}
	}

	tree.Remove([]byte("aardvark"))
	want := len(words) - 1
	if got := tree.Len(); got != want {
		t.Errorf("Len after remove\n got: %d\nwant: %d", got, want)
	}

	// Insert on an existing key.
	tree.Insert([]byte("toad"), 5)
	if got := tree.Len(); got != want {
		t.Errorf("Len after insert with existing key\n got: %d\nwant: %d", got, want)
	}
}

func TestLongestPrefix(t *testing.T) {
	if got, ok := New[int]().LongestPrefix([]byte("a")); ok || got != 0 {
		t.Errorf("LongestPrefix on empty tree\n got: (%v, %t)\nwant: (0, false)", got, ok)
	}

	tree := build(words)
	key := []byte("winkley")
	want := "winkle"
	if got, ok := tree.LongestPrefix(key); !ok || got != want {
		t.Errorf("LongestPrefix(%s)\n got: (%s, %t)\nwant: (%s, true)", key, got, ok, want)
	}

	want = "wink"
	if got, ok := tree.LongestPrefix([]byte(want)); !ok || got != want {
		t.Errorf("LongestPrefix(%s)\n got: (%s, %t)\nwant: (%s, true)", key, got, ok, want)
	}
}

func TestMax(t *testing.T) {
	if got, ok := New[int]().Max(); ok || got != 0 {
		t.Errorf("Max on empty tree\n got: (%v, %t)\nwant: (0, false)", got, ok)
	}

	tree := build(words)

	want := words[len(words)-1]
	if got, ok := tree.Max(); !ok || got != want {
		t.Errorf("Max\n got: (%s, %t)\nwant: (%s, true)", got, ok, want)
	}

	want = "zzz"
	tree.Insert([]byte(want), want)
	if got, ok := tree.Max(); !ok || got != want {
		t.Errorf("Max after insert\n got: (%s, %t)\nwant: (%s, true)", got, ok, want)
	}
}

func TestMin(t *testing.T) {
	if got, ok := New[int]().Min(); ok || got != 0 {
		t.Errorf("Min on empty tree\ngot: (%v, %t)\nwant: (0, false)", got, ok)
	}

	tree := build(words)
	want := words[0]
	if got, ok := tree.Min(); !ok || got != want {
		t.Errorf("Min\n got: (%s, %v)\nwant: (%s, true)", got, ok, want)
	}

	want = "a"
	tree.Insert([]byte(want), want)
	if got, ok := tree.Min(); !ok || got != want {
		t.Errorf("Min after insert\n got: (%s, %t)\nwant: (%s, true)", got, ok, want)
	}
}

func TestPredecessor(t *testing.T) {
	if got, ok := New[int]().Predecessor([]byte("key")); ok || got != 0 {
		t.Errorf("Predecessor on empty tree\n got: (%v, %t)\nwant: (0, false)", got, ok)
	}

	tree := build(words)

	firstKey := []byte(words[0])
	if got, ok := tree.Predecessor(firstKey); ok || got != "" {
		t.Errorf("Predecessor with the first key\n got: (%v, %t)\nwant: (\"\", false)", got, ok)
	}

	key := []byte("non-existent key")
	if got, ok := tree.Predecessor(key); ok || got != "" {
		t.Errorf("Predecessor(%s)\n got: (%s, %t)\nwant: (\"\", false)", key, got, ok)
	}

	for i, name := range words[1:] {
		want := words[i]
		if got, ok := tree.Predecessor([]byte(name)); !ok || got != want {
			t.Errorf("Predecessor(%s)\n got: (%s, %t)\nwant: (%s, true)", name, got, ok, want)
		}
	}
}

func TestRemove(t *testing.T) {
	tree := build(words)

	if got, ok := tree.Remove([]byte("aardvs")); ok || got != "" {
		t.Errorf("Remove with a key that doesn't exist\n got: (%s, %t)\nwant: (\"\", false)", got, ok)
	}

	for _, want := range words {
		if got, ok := tree.Remove([]byte(want)); !ok || got != want {
			t.Errorf("Remove(%s)\n got: (%s, %t)\nwant: (%s, true)", want, got, ok, want)
		}
	}

	if got, ok := tree.Remove([]byte{0}); ok || got != "" {
		t.Errorf("Remove on empty tree\n got: (%s, %t)\nwant: (\"\", false)", got, ok)
	}

	tree = build(words)

	// Attempt to delete an existing node that doesn't have a value.
	if got, ok := tree.Remove([]byte("aard")); ok || got != "" {
		t.Errorf("Remove node that doesn't have a value\n got: (%v, %t)\nwant: (\"\", false)", got, ok)
	}

	// Removing wit will cause a merge with the parent (wi) and the parent's
	// only remaining child (ll).
	want := "wit"
	if got, ok := tree.Remove([]byte(want)); !ok || got != want {
		t.Errorf("Remove(%s)\n got: (%s, %t)\nwant: (%s, true)", want, got, ok, want)
	}

	// Attempt to delete a key that doesn't exist but is a prefix of another
	// existing key.
	if got, ok := tree.Remove([]byte("ba")); ok || got != "" {
		t.Errorf("Remove node that doesn't have a value\n got: (%v, %t)\nwant: (\"\", false)", got, ok)
	}
}

func TestSuccessor(t *testing.T) {
	if got, ok := New[int]().Successor([]byte("key")); ok || got != 0 {
		t.Errorf("Successor on empty tree\n got: (%v, %t)\nwant: (0, false)", got, ok)
	}

	tree := build(words)

	last := words[len(words)-1]
	if got, ok := tree.Successor([]byte(last)); ok || got != "" {
		t.Errorf("Successor with the last key\n got: (%v, %t)\nwant: (\"\", false)", got, ok)
	}

	key := []byte("non-existent key")
	if got, ok := tree.Successor(key); ok || got != "" {
		t.Errorf("Successor(%s)\n got: (%s, %t)\nwant: (\"\", false)", key, got, ok)
	}

	for i, name := range words[:len(words)-1] {
		want := words[i+1]
		if got, ok := tree.Successor([]byte(name)); !ok || got != want {
			t.Errorf("Successor(%s)\n got: (%s, %t)\nwant: (%s, true)", name, got, ok, want)
		}
	}
}

func TestValues(t *testing.T) {
	want := make([]string, 0, len(words))
	if got := New[int]().Values(); len(got) != 0 {
		t.Errorf("Values returned non-empty slice for empty tree\n got: %v\nwant: %v", got, want)
	}

	tree := build(words)
	for _, word := range words {
		want = append(want, word)
	}
	if got := tree.Values(); !reflect.DeepEqual(got, want) {
		t.Errorf("Values\n got: %v\nwant: %v", got, want)
	}
}

func TestWalk(t *testing.T) {
	tree := build(words)
	limit := 3
	got := make([]string, 0, limit)
	i := 0
	// Walk the first 3 words rooted at "to".
	tree.Walk([]byte("to"), func(value string) bool {
		got = append(got, value)
		i++
		return i < limit
	})

	want := hasPrefix("to", words)[:limit]
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Walk\n got: %v\nwant: %v", got, want)
	}

	// Walk the entire tree.
	values := make([]string, 0, len(words))
	tree.Walk([]byte{}, func(value string) bool {
		values = append(values, value)
		return true
	})
	if !reflect.DeepEqual(values, words) {
		t.Errorf("Walk\n got: %v\nwant: %v", values, words)
	}
}

var words = []string{
	"aardvark",
	"aardwolf",
	"abacus",
	"babble",
	"backtrack",
	"beehive",
	"create",
	"macro",
	"macroanalysis",
	"macroanalyst",
	"macrochelys",
	"mactroid",
	"obsequious",
	"sequence",
	"to",
	"toa",
	"toad",
	"toady",
	"toadyism",
	"what",
	"win",
	"wink",
	"winkle",
	"winkleman",
	"will",
	"wilting",
	"wit",
}
