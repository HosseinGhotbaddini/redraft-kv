package test

import (
	"os"
	"testing"

	"github.com/HosseinGhotbaddini/redraft-kv/store"
)

func TestBoltStore_SetGetDelete(t *testing.T) {
	tempFile := "testdata/bolt-test.db"
	_ = os.MkdirAll("testdata", 0755)
	defer os.Remove(tempFile)

	kv := store.NewBoltStore(tempFile)
	defer kv.Close()

	t.Run("set and get", func(t *testing.T) {
		if err := kv.Set("bolt-key", []byte("value")); err != nil {
			t.Fatal("set failed:", err)
		}

		val, err := kv.Get("bolt-key")
		if err != nil {
			t.Fatal("get failed:", err)
		}
		if string(val) != "value" {
			t.Errorf("expected 'value', got '%s'", val)
		}
	})

	t.Run("overwrite key", func(t *testing.T) {
		kv.Set("bolt-key", []byte("new"))
		val, _ := kv.Get("bolt-key")
		if string(val) != "new" {
			t.Errorf("expected 'new', got '%s'", val)
		}
	})

	t.Run("delete key", func(t *testing.T) {
		kv.Delete("bolt-key")
		val, _ := kv.Get("bolt-key")
		if val != nil {
			t.Error("expected key to be deleted")
		}
	})
}
