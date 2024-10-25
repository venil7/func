package function

type Lazy[A any] func() A
type Mapping[A any, B any] func(a A) B
type Predicate[A any] Mapping[A, bool]
type LazyErr[A any] func() (A, error)
type MapLazyErr[A any, B any] func(a A) (B, error)
