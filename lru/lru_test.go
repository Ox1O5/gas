package lru

import (
	"testing"
)

type String string

func (s String)Len() int{
	return len(s)
}

var getTests = []struct{
	name string
	keyToAdd interface{}
	keyToGet interface{}
	expectedOK bool
}{
	{"string_hit", "key1", "key1", true},
	{"string_miss", "key2", "notKey", false},
}

func TestCache_Add(t *testing.T) {
	c := New(int64(0), nil)
	var expectedNbyte int64 = 0
	for _, tt :=  range getTests {
		c.Add(tt.keyToAdd.(string), String("1234"))
		expectedNbyte += int64(len(tt.keyToAdd.(string)) + String("1234").Len())
	}

	if c.nbytes != expectedNbyte {
		t.Fatalf("expected %v, but got %v", expectedNbyte, c.nbytes)
	}
	if c.Len() != len(getTests) {
		t.Fatalf("expected len: %v, but got %v", len(getTests), c.Len())
	}
}

func TestCache_Get(t *testing.T) {
	c := New(int64(0), nil)
	for _, tt := range getTests {
		c.Add(tt.keyToAdd.(string), String("1234"))
		val, ok := c.Get(tt.keyToGet.(string))
		if ok != tt.expectedOK {
			t.Fatalf("%s: cache hit = %v; want %v", tt.name, ok, !ok)
		} else if ok && string(val.(String)) != "1234" {
			t.Fatalf("%s expected get to return \"1234\" but got %v", tt.name, val)
		}
	}
}



func TestCache_RemoveOldest(t *testing.T) {
	var nbyte int64 = 0
	for _, tt :=  range getTests {
		nbyte += int64(len(tt.keyToAdd.(string)) + String("1234").Len())
	}

	c := New(nbyte-1, nil)
	for _, tt := range getTests {
		c.Add(tt.keyToAdd.(string), String("1234"))
	}

	if _, hitKey1 := c.Get(getTests[0].keyToAdd.(string)); hitKey1 || c.Len() != 1 {
		t.Fatalf("RemoveOldsest key1 failed")
	}
}


