package storage

import (
	"strconv"
	"testing"
)

type testCase struct {
	name  string
	key   string
	value string
	kind  Kind
}

func TestSetGetWithKind(t *testing.T) {
	cases := []testCase{
		{"hello world", "hello", "world", KindString},
		{"YYYXXX 123", "YYYXXX", "123", KindInt},
		{"U1", "U1", "", KindUndefined},
		{"U2 undefined", "U2", "undefined", KindUndefined},
	}

	s, err := NewStorage()
	if err != nil {
		t.Errorf("new storage: %v", err)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s.Set(c.key, c.value)
			sValue := s.Get(c.key)

			if *sValue != c.value {
				t.Errorf("values not equal")
			}

			if getType(*sValue) != getType(c.value) {
				t.Errorf("value kinds not equal")
			}

			if getType(*sValue) != c.kind {
				t.Errorf("expected value kind: %v", c.kind)
			}
		})
	}
}

func BenchmarkStorage_Set(b *testing.B) {
	storage, err := NewStorage()
	if err != nil {
		b.Fatalf("Error creating storage: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}
}

func BenchmarkStorage_Get(b *testing.B) {
	storage, err := NewStorage()
	if err != nil {
		b.Fatalf("Error creating storage: %v", err)
	}

	for i := 0; i < 1000; i++ {
		storage.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		storage.Get("key" + strconv.Itoa(i%1000))
	}
}

func BenchmarkStorage_SetThenGet(b *testing.B) {

	storage, err := NewStorage()
	if err != nil {
		b.Fatalf("Error creating storage: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "key" + strconv.Itoa(i)
		storage.Set(key, "value"+strconv.Itoa(i))
		storage.Get(key)
	}
}
