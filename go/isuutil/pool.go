package isuutil

import "sync"

var ()

type ArrPool[T any] struct {
	data *sync.Pool
}

func NewArrPool[T any](defaultSize int) *ArrPool[T] {
	return &ArrPool[T]{
		data: &sync.Pool{
			New: func() interface{} {
				s := make([]T, 0, defaultSize)
				return &s
			},
		},
	}
}

func (p *ArrPool[T]) get() ([]T, func()) {
	ptr := p.data.Get().(*[]T)
	arr := *ptr
	return arr, func() {
		arr = arr[0:0]
		*ptr = arr
		p.data.Put(ptr)
	}
}

type Pool[T any] struct {
	data *sync.Pool
	putF func(T) T
}

func NewPool[T any](fn func() T, putF func(T) T) *Pool[T] {
	return &Pool[T]{
		data: &sync.Pool{
			New: func() interface{} {
				return fn()
			},
		},
		putF: putF,
	}
}

func (p *Pool[T]) get() (T, func()) {
	res := p.data.Get().(T)
	return res, func() {
		p.data.Put(p.putF(res))
	}
}
