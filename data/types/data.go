package data

type Numeric[T any] interface {
	FromInt64(int64) T
	FromUint64(uint64) T
	Zero() T
	One() T
	
	Add(T) T
	Sub(T) T
	Mul(T) T
	Div(T) T
	
	LT(T) bool
	EQ(T) bool
}

var _ Numeric[Uint256] = Uint256{}

func (x Uint256) Add(y Uint256) Uint256 {
	return Uint256{
		A: x.A+y.A,
		// ...
	}
}