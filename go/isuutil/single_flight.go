package isuutil

import "golang.org/x/sync/singleflight"

// SingleFlightGroup は singleflight.Group を型パラメーターに対応して使いやすくしたもの。
type SingleFlightGroup[T any] struct {
	group *singleflight.Group
}

func NewSingleFlightGroup[T any]() *SingleFlightGroup[T] {
	return &SingleFlightGroup[T]{group: &singleflight.Group{}}
}

func (g *SingleFlightGroup[T]) Do(key string, fn func() (T, error)) (T, error, bool) {
	got, err, shared := g.group.Do(key, func() (interface{}, error) {
		return fn()
	})

	v, ok := got.(T)
	if ok {
		return v, err, shared
	}
	var zero T
	return zero, err, shared
}
