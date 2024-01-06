package pokecache

import (
	"testing"
	"time"
	"fmt"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second

	cases := []struct{
		key string
		value []byte
	}{
		{
			key: "case1",
			value: []byte("case1_value"),
		},
		{
			key: "case2",
			value: []byte("case2_value"),
		},
	}

	for i, c := range cases{
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.value)
			val, ok := cache.Get(c.key)

			if !ok {
				t.Errorf("Key %v not found", c.key)
				return
			}

			if string(val) != string(c.value) {
				t.Errorf("Value of key %v is not correct", c.key)
				return
			}
		})
	}

}

func TestReapLoop(t *testing.T){
	const baseTime = 2 * time.Second
	const waitTime = baseTime + 5 * time.Second

	cache := NewCache(baseTime)

	key := "https://testcase1.com"
	value := []byte("value_testcase1")

	cache.Add(key, value)

	_, ok := cache.Get(key)

	if !ok {
		t.Errorf("Key %v not found", key)
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get(key)

	if ok {
		t.Errorf("Key %v should not be present.", key)
		return
	}

}
