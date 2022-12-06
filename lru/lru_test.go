package lru

import (
	"reflect"
	"testing"
)

// 自定义类型
type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("name", String("sean"))
	// 测试一个能获取到例子
	if ele, ok := lru.Get("name"); !ok || ele.(String) != "sean" {
		t.Fatal("key exist but not correct")
	}
	if _, ok := lru.Get("key"); ok {
		t.Fatal("key not exist but found")
	}
}

func TestRemoveOldest(t *testing.T) {

	// 测试存储2个元素,第三个元素会超过最大值清除
	lru := New(int64(18), nil)
	lru.Add("name1", String("sean"))
	lru.Add("name2", String("emma"))
	lru.Add("name3", String("alex"))

	// 检查name1存在
	if _, ok := lru.Get("name"); ok || lru.Len() != 2 {
		t.Fatal("key removeOldest failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("v2"))
	lru.Add("k3", String("v3"))
	lru.Add("k4", String("v4"))

	// 通过触发空间溢出触发回调收集删除的key
	expected := []string{"key1", "k2"}

	if !reflect.DeepEqual(keys, expected) {
		t.Fatal("callback failed")
	}

}
