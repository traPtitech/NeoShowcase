package optional

type Of[T any] struct {
	V     T
	Valid bool
}

func New[T any](v T, valid bool) Of[T] {
	return Of[T]{
		V:     v,
		Valid: valid,
	}
}

func From[T any](v T) Of[T] {
	return Of[T]{
		V:     v,
		Valid: true,
	}
}

// ValueOrZero 値が入っているときはその値を、そうでないときはゼロ値を返します。
func (o Of[T]) ValueOrZero() T {
	if o.Valid {
		return o.V
	}
	var t T
	return t
}
