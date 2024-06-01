package store

import (
	"fmt"
	"sync"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	t.Run("Simple Set and Get", func(t *testing.T) {
		Set("key1", "value1")
		got := Get("key1")
		want := "value1"
		if got != want {
			t.Errorf("Get(key1) = %q; want %q", got, want)
		}
	})

	t.Run("Update Value", func(t *testing.T) {
		Set("key1", "newValue")
		got := Get("key1")
		want := "newValue"
		if got != want {
			t.Errorf("Get(key1) = %q; want %q", got, want)
		}
	})
}

func TestGetNonExistentKey(t *testing.T) {
	t.Run("Get Non-Existent Key", func(t *testing.T) {
		got := Get("nonexistent")
		want := ""
		if got != want {
			t.Errorf("Get(nonexistent) = %q; want %q", got, want)
		}
	})
}

func TestStrln(t *testing.T) {
	t.Run("Length of Existing Key", func(t *testing.T) {
		Set("key2", "value2")
		got := Strln("key2")
		want := 6
		if got != want {
			t.Errorf("Strln(key2) = %d; want %d", got, want)
		}
	})

	t.Run("Length of Non-Existent Key", func(t *testing.T) {
		got := Strln("nonexistent")
		want := 0
		if got != want {
			t.Errorf("Strln(nonexistent) = %d; want %d", got, want)
		}
	})
}

func TestDel(t *testing.T) {
	t.Run("Delete Existing Key", func(t *testing.T) {
		Set("key3", "value3")
		got := Del("key3")
		want := "value3"
		if got != want {
			t.Errorf("Del(key3) = %q; want %q", got, want)
		}

		got = Get("key3")
		want = ""
		if got != want {
			t.Errorf("Get(key3) after Del = %q; want %q", got, want)
		}
	})

	t.Run("Delete Non-Existent Key", func(t *testing.T) {
		got := Del("nonexistent")
		want := ""
		if got != want {
			t.Errorf("Del(nonexistent) = %q; want %q", got, want)
		}
	})
}

func TestAppend(t *testing.T) {
	t.Run("Append to Existing Key", func(t *testing.T) {
		Set("key4", "value4")
		Append("key4", "appended")
		got := Get("key4")
		want := "value4appended"
		if got != want {
			t.Errorf("Get(key4) after Append = %q; want %q", got, want)
		}
	})

	t.Run("Append to Non-Existent Key", func(t *testing.T) {
		Append("nonexistent", "appended")
		got := Get("nonexistent")
		want := "appended"
		if got != want {
			t.Errorf("Get(nonexistent) after Append = %q; want %q", got, want)
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	var wg sync.WaitGroup
	key := "concurrentKey"

	t.Run("Concurrent Set", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				Set(key, fmt.Sprintf("value%d", i))
			}(i)
		}

		wg.Wait()
		if Strln(key) == 0 {
			t.Errorf("Strln(concurrentKey) should not be 0 after concurrent Set operations")
		}
	})

	t.Run("Concurrent Append", func(t *testing.T) {
		Set(key, "")
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				Append(key, fmt.Sprintf("%d", i))
			}(i)
		}

		wg.Wait()
		if Strln(key) == 0 {
			t.Errorf("Strln(concurrentKey) should not be 0 after concurrent Append operations")
		}
	})
}
