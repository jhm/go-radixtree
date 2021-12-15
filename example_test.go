package radixtree

import "fmt"

func ExampleRadixTree_Contains() {
	t := New[int]()
	t.Insert([]byte("John"), 1)
	fmt.Println(t.Contains([]byte("John")))
	fmt.Println(t.Contains([]byte("Bill")))
	// Output:
	// true
	// false
}

func ExampleRadixTree_Find() {
	t := New[int]()
	t.Insert([]byte("John"), 1)
	t.Insert([]byte("Jonathan"), 2)

	vs := t.Find([]byte("Jo"))
	fmt.Println(vs)
	// Output:
	// [1 2]
}

func ExampleRadixTree_Get() {
	t := New[int]()
	t.Insert([]byte("John"), 1)

	v, found := t.Get([]byte("John"))
	fmt.Println(found)
	fmt.Println(v)
	// Output:
	// true
	// 1
}

func ExampleRadixTree_Insert() {
	t := New[int]()
	old, found := t.Insert([]byte("John"), 1)
	fmt.Println(old)
	fmt.Println(found)

	old, found = t.Insert([]byte("John"), 2)
	fmt.Println(old)
	fmt.Println(found)
	// Output:
	// 0
	// false
	// 1
	// true
}

func ExampleRadixTree_LongestPrefix() {
	t := New[int]()
	t.Insert([]byte("Eric"), 1)
	v, found := t.LongestPrefix([]byte("Ericson"))
	fmt.Println(v)
	fmt.Println(found)
	// Output:
	// 1
	// true
}

func ExampleRadixTree_Max() {
	t := New[int]()
	t.Insert([]byte("Aaron"), 1)
	t.Insert([]byte("Zaire"), 2)
	v, found := t.Max()
	fmt.Println(v)
	fmt.Println(found)
	// Output:
	// 2
	// true
}

func ExampleRadixTree_Min() {
	t := New[int]()
	t.Insert([]byte("Aaron"), 1)
	t.Insert([]byte("Zaire"), 2)
	v, found := t.Min()
	fmt.Println(v)
	fmt.Println(found)
	// Output:
	// 1
	// true
}

func ExampleRadixTree_Predecessor() {
	t := New[int]()
	t.Insert([]byte("Aaron"), 1)
	t.Insert([]byte("Zaire"), 2)
	v, found := t.Predecessor([]byte("Zaire"))
	fmt.Println(v)
	fmt.Println(found)
	// Output:
	// 1
	// true
}

func ExampleRadixTree_Remove() {
	t := New[int]()
	t.Insert([]byte("Aaron"), 1)
	v, found := t.Remove([]byte("Aaron"))
	fmt.Println(v)
	fmt.Println(found)
	fmt.Println(t.Len())
	// Output:
	// 1
	// true
	// 0
}

func ExampleRadixTree_Successor() {
	t := New[int]()
	t.Insert([]byte("Aaron"), 1)
	t.Insert([]byte("Zaire"), 2)
	v, found := t.Successor([]byte("Aaron"))
	fmt.Println(v)
	fmt.Println(found)
	// Output:
	// 2
	// true
}

func ExampleRadixTree_Values() {
	t := New[int]()
	t.Insert([]byte("Zaire"), 0)
	t.Insert([]byte("Aaron"), 1)
	t.Insert([]byte("Erica"), 2)
	fmt.Println(t.Values())
	// Output:
	// [1 2 0]
}
