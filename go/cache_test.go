package main

import (
	stdcmp "cmp"
	"github.com/google/go-cmp/cmp"
	"slices"
	"testing"
)

func TestCache_Values(t *testing.T) {
	c := NewCache[string, ExampleItem]()
	c.Set("name1", ExampleItem{Name: "name1"})
	c.Set("name2", ExampleItem{Name: "name2"})

	want := []ExampleItem{
		{Name: "name1"},
		{Name: "name2"},
	}

	got := c.Values()

	slices.SortFunc(want, func(a, b ExampleItem) int {
		return stdcmp.Compare(a.Name, b.Name)
	})
	slices.SortFunc(got, func(a, b ExampleItem) int {
		return stdcmp.Compare(a.Name, b.Name)
	})

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

}

func TestCache_Get(t *testing.T) {
	c := NewCache[int, ExampleItem]()
	c.Set(1, ExampleItem{Name: "name1"})
	c.Set(2, ExampleItem{Name: "name2"})

	want := ExampleItem{Name: "name2"}
	got, ok := c.Get(2)
	if !ok {
		t.Errorf("not found")
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

	want = ExampleItem{}
	got, ok = c.Get(3)
	if ok {
		t.Errorf("found")
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

}

func TestCache_Delete(t *testing.T) {
	c := NewCache[string, ExampleItem]()

	c.Set("name1", ExampleItem{Name: "name1"})
	c.Set("name2", ExampleItem{Name: "name2"})

	c.Delete("name1")

	_, ok := c.Get("name1")
	if ok {
		t.Errorf("found")
	}

}

func TestSliceCache_Get(t *testing.T) {
	sc := NewSliceCache[string, ExampleItem]()
	sc.Append("key", ExampleItem{Name: "name1"})
	sc.Append("key", ExampleItem{Name: "name2"})

	want := []ExampleItem{
		{Name: "name1"},
		{Name: "name2"},
	}

	got := sc.Get("key")

	slices.SortFunc(want, func(a, b ExampleItem) int {
		return stdcmp.Compare(a.Name, b.Name)
	})
	slices.SortFunc(got, func(a, b ExampleItem) int {
		return stdcmp.Compare(a.Name, b.Name)
	})

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func TestSliceCache_Flush(t *testing.T) {
	sc := NewSliceCache[string, ExampleItem]()
	sc.Append("key", ExampleItem{Name: "name1"})
	sc.Append("key", ExampleItem{Name: "name2"})

	sc.Flush()

	got := sc.Get("key")
	if len(got) != 0 {
		t.Errorf("not flushed")
	}
}
