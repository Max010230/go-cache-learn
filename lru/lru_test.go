package lru

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestCache_Get(t *testing.T) {
	callback := func(key string, value Value) {
		result, _ := json.Marshal(value)
		log.Println(fmt.Sprintf("%s:%s", key, string(result)))
	}
	cache := New(int64(1024*100), callback)
	for i := 0; i < 10000; i++ {
		cache.Add("key"+fmt.Sprint(i), String("Hello Word!"))
		cache.Get("key1")
	}
	if v, ok := cache.Get("key1"); !ok || string(v.(String)) != "Hello Word!" {
		t.Fatalf("cache hit failed")
	}
	t.Log(cache.Len())
}
