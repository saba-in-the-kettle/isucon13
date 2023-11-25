package main

import (
	"context"
	"github.com/patrickmn/go-cache"
	"strconv"
	"sync"
	"time"
)

var (
	exampleNameCache = NewCache[string, ExampleItem]()
	examplesCache    = NewSliceCache[string, ExampleItem]()
)

// InitCache はmain関数とinitializeHandlerの両方で呼び出す必要がある
func InitCache(ctx context.Context) error {
	exampleNameCache.Flush()

	e := ExampleItem{Name: "example"}
	exampleNameCache.Set(e.Name, e)

	examplesCache.Append("map key", e)

	return nil
}

// Cache は Type Parameterに対応したKey Value Storeです
type Cache[K comparable, V any] struct {
	cache *cache.Cache
}

// NewCache はCacheを初期化します。
// K はキーの型、V は値の型です。
func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{cache: cache.New(cache.NoExpiration, cache.NoExpiration)}
}

// NewCacheWithExpire は有効期限付きのCacheを初期化します。
// K はキーの型、V は値の型です。
func NewCacheWithExpire[K comparable, V any](defaultExpiration, cleanupInterval time.Duration) *Cache[K, V] {
	return &Cache[K, V]{cache: cache.New(defaultExpiration, cleanupInterval)}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	v, ok := c.cache.Get(toKey(key))
	if ok {
		return v.(V), true
	}
	var defaultValue V
	return defaultValue, false
}

func (c *Cache[K, V]) Set(k K, v V) {
	c.cache.Set(toKey(k), v, cache.DefaultExpiration)
}

func (c *Cache[K, V]) Delete(k K) {
	c.cache.Delete(toKey(k))
}

// Values はキャッシュの値を全て取得します
// 並び順は順不同です。
func (c *Cache[K, V]) Values() []V {
	var values []V
	for _, v := range c.cache.Items() {
		values = append(values, v.Object.(V))
	}
	return values
}

func (c *Cache[K, V]) SetWithExpire(k K, v V, d time.Duration) {
	c.cache.Set(toKey(k), v, d)
}

// Flush はキャッシュをクリアします
func (c *Cache[K, V]) Flush() {
	c.cache.Flush()
}

// SliceCache は []V をキャッシュする構造体です
type SliceCache[K comparable, V any] struct {
	item map[K][]V
	sync.RWMutex
}

func NewSliceCache[K comparable, V any]() *SliceCache[K, V] {
	return &SliceCache[K, V]{
		item:    map[K][]V{},
		RWMutex: sync.RWMutex{},
	}
}

func (sc *SliceCache[K, V]) Get(key K) []V {
	sc.RLock()
	defer sc.RUnlock()

	return sc.item[key]
}

func (sc *SliceCache[K, V]) Append(key K, value V) {
	sc.Lock()
	defer sc.Unlock()

	if len(sc.item[key]) == 0 {
		sc.item[key] = []V{}
	}

	sc.item[key] = append(sc.item[key], value)
}

func (sc *SliceCache[K, V]) Flush() {
	sc.Lock()
	defer sc.Unlock()

	sc.item = map[K][]V{}
}

func toKey[K comparable](key K) string {
	a := any(key)
	switch a.(type) {
	case int, int8, int16, int32, int64:
		return strconv.Itoa(a.(int))
	case uint, uint8, uint16, uint32, uint64:
		return strconv.Itoa(int(a.(uint)))
	case float32:
		return strconv.FormatFloat(float64(a.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(a.(float64), 'f', -1, 64)
	case string:
		return a.(string)
	default:
		panic("not supported type")
	}
}

type ExampleItem struct {
	Name string
}
