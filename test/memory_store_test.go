package test

import (
	"testing"

	"github.com/HosseinGhotbaddini/redraft-kv/store"
)

func TestStore_SetGet(t *testing.T) {
	kv := store.NewMemoryStore()

	t.Run("set and get basic key", func(t *testing.T) {
		kv.Set("foo", []byte("bar"))
		val, err := kv.Get("foo")
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		if val == nil {
			t.Fatal("expected key 'foo' to exist")
		}
		if string(val) != "bar" {
			t.Errorf("expected 'bar', got '%s'", val)
		}
	})

	t.Run("overwrite existing key", func(t *testing.T) {
		kv.Set("foo", []byte("baz"))
		val, _ := kv.Get("foo")
		if string(val) != "baz" {
			t.Errorf("expected 'baz', got '%s'", val)
		}
	})

	t.Run("get non-existent key", func(t *testing.T) {
		val, _ := kv.Get("unknown")
		if val != nil {
			t.Error("expected key 'unknown' to not exist")
		}
	})

	t.Run("set and get empty key and value", func(t *testing.T) {
		kv.Set("", []byte{})
		val, _ := kv.Get("")
		if val == nil || len(val) != 0 {
			t.Errorf("expected empty value, got %v", val)
		}
	})
}

func TestStore_Delete(t *testing.T) {
	kv := store.NewMemoryStore()

	t.Run("delete existing key", func(t *testing.T) {
		kv.Set("delete-me", []byte("temp"))
		kv.Delete("delete-me")
		val, _ := kv.Get("delete-me")
		if val != nil {
			t.Error("expected key 'delete-me' to be deleted")
		}
	})

	t.Run("delete non-existent key (should not panic)", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("expected safe delete, got panic: %v", r)
			}
		}()
		kv.Delete("ghost")
	})
}
