package optional

type Of[T any] struct {
	V     T
	Valid bool
}

func None[T any]() Of[T] {
	return Of[T]{
		Valid: false,
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

func Map[T, U any](o Of[T], mapper func(T) U) Of[U] {
	if o.Valid {
		return From(mapper(o.V))
	}
	return None[U]()
}

func MapErr[T, U any](o Of[T], mapper func(T) (U, error)) (Of[U], error) {
	if o.Valid {
		v, err := mapper(o.V)
		return From(v), err
	}
	return None[U](), nil
}

func FromPtr[T any](v *T) Of[T] {
	if v == nil {
		return None[T]()
	}
	return From(*v)
}

func FromNonZero[T comparable](v T) Of[T] {
	var zero T
	if v == zero {
		return None[T]()
	}
	return From(v)
}
