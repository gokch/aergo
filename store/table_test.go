package store

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/aergoio/aergo-lib/db"
)

func TestTableDatabase(t *testing.T) {
	testTableDatabase(t, []byte("prefix"))
}

func TestEmptyPrefixTableDatabase(t *testing.T) {
	testTableDatabase(t, []byte(""))
}

func testTableDatabase(t *testing.T, prefix []byte) {
	table := NewTable(db.NewDB(db.MemoryImpl, ""), prefix)

	var entries = []struct {
		key   []byte
		value []byte
	}{
		{[]byte{0x01, 0x02}, []byte{0x0a, 0x0b}},
		{[]byte{0x03, 0x04}, []byte{0x0c, 0x0d}},
		{[]byte{0x05, 0x06}, []byte{0x0e, 0x0f}},

		{[]byte{0xff, 0xff, 0x01}, []byte{0x1a, 0x1b}},
		{[]byte{0xff, 0xff, 0x02}, []byte{0x1c, 0x1d}},
		{[]byte{0xff, 0xff, 0x03}, []byte{0x1e, 0x1f}},
	}

	for _, entry := range entries {
		table.Set(entry.key, entry.value)
	}
	for _, entry := range entries {
		got := table.Get(entry.key)
		if !bytes.Equal(got, entry.value) {
			t.Fatalf("Value mismatch: want=%v, got=%v", entry.value, got)
		}
	}

	// Test batch operation
	table = NewTable(db.NewDB(db.MemoryImpl, ""), prefix)
	batch := table.NewBulk()
	for _, entry := range entries {
		batch.Set(entry.key, entry.value)
	}
	batch.Flush()
	for _, entry := range entries {
		got := table.Get(entry.key)
		if !bytes.Equal(got, entry.value) {
			t.Fatalf("Value mismatch: want=%v, got=%v", entry.value, got)
		}
	}

	check := func(iter db.Iterator, expCount, index int) {
		count := 0
		for iter.Valid() {
			key, value := iter.Key(), iter.Value()
			fmt.Println(key, value)
			if !bytes.Equal(key, entries[index].key) {
				t.Fatalf("Key mismatch: want=%v, got=%v", entries[index].key, key)
			}
			if !bytes.Equal(value, entries[index].value) {
				t.Fatalf("Value mismatch: want=%v, got=%v", entries[index].value, value)
			}
			index += 1
			count++
			iter.Next()
		}
		if count != expCount {
			t.Fatalf("Wrong number of elems, exp %d got %d", expCount, count)
		}
	}
	// Test iterators
	// check(table.Iterator(nil, nil), 6, 0)
	// Test iterators with prefix
	check(table.Iterator([]byte{0xff, 0xff}, nil), 3, 3)
	// Test iterators with start point
	check(table.Iterator(nil, []byte{0xff, 0xff, 0x02}), 2, 4)
	// Test iterators with prefix and start point
	check(table.Iterator([]byte{0xee}, nil), 0, 0)
	check(table.Iterator(nil, []byte{0x00}), 6, 0)
}
