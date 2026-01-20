package network

type dto[T any] struct {
	value *T
}

func (d *dto[T]) GetValue() *T {
	return d.value
}

func NewDto[T any](value *T) Dto[T] {
	return &dto[T]{value: value}
}
