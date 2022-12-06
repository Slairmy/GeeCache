package lru

import "container/list"

type Cache struct {
	maxBytes int64 // 允许使用的最大内存
	nbytes   int64 // 当前已经使用的内存
	ll       *list.List
	cache    map[string]*list.Element
	//
	OnEvicted func(key string, value Value) // 定义函数,当元素淘汰的时候执行
}

// entry 链表中存储的值(值是键值对) -- 保存key的原因是需要再淘汰链表节点的时候同时可以删除 map 中的数据
type entry struct {
	key   string
	value Value
}

type Value interface { // 缓存中的任意值类型
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// 被使用了需要移动到队列尾部
		c.ll.MoveToFront(ele)

		kv := ele.Value.(*entry)
		return kv.value, true
	}

	return
}

// RemoveOldest 删除元素 -- 获取队列头部节点进行删除
func (c *Cache) RemoveOldest() {
	// 取出链表头部的节点
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		// 同时删除map中的数据
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		// 减少当前使用的内存大小
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 新增/修改元素
func (c *Cache) Add(key string, value Value) {
	// 新增元素,应该是需要判断当前有没有满如果满了需要进行删除
	if ele, ok := c.cache[key]; ok {
		// 修改元素
		// 移动到末尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)

		// 改变当前暂用大小,使用修改后的大小减去修改前的大小
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 新增
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele

		// 增加当前内存使用
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// 当新增一个元素或者修改一个元素之后超出了最大使用空间,直接移除使用最少的元素
	if c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return len(c.cache)
}
