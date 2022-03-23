package operation

import "github.com/mconcat/dbci/operation/types"

type UnaryFunction[T, U any] interface {
	IsUnaryFunction(T) U
	ProtoFunction() types.Function
}

type Numeric interface {
	PrimitiveNumeric | ~Int256 | ~Uint256
}

type PrimitiveNumeric interface {
	~int64 | ~uint64
}

type FunctionArithmetic[T Numeric] struct { types.FunctionNumeric }
var _ UnaryFunction[T, T] = FunctionArithmetic[T]
func (FunctionArithmetic[T]) IsUnaryFunction(_ T) (res T) { return }
func (fn FunctionArithmetic[T]) ProtoFunction() types.Function { return fn.FunctionNumeric }
func numeric(op types.FunctionNumeric_Operator, v Numeric) types.FunctionNumeric {
	var operand isFunctionNumeric_Operand
	switch v := v.(type) {
	case int64:
		operand = FunctionNumeric_Int64Operand{ Int64Operand: v }
	case uint64:
		operand = FunctionNumeric_Uint64Operand{ Uint64Operand: v }
	case Int256:
		operand = FunctionNumeric_Int256Operand{ Int256Operand: v }
	case Uint256:
		operand = FunctionNumeric_Uint256Operand{ Uint256Operand: v }
	default:
		panic("fdsafsa")
	}
	return types.FunctionNumeric {
		Operator: op,
		Operand: operand,
	}
}

func Add[T Numeric](v T) FunctionArithmetic[T] {
	return FunctionArithmetic[T]{ numeric(types.FunctionNumeric_ADD, v) }
}
func Sub[T Numeric](v T) FunctionArithmetic[T] {
	return FunctionArithmetic[T]{ numeric(types.FunctionNumeric_Sub, v) }
}

type FunctionNumericComparative[T Numeric] struct { types.FunctionNumeric }
var _ UnaryFunction[T, bool] = FunctionNumericComparative[T]
func (FunctionNumericComparative[T]) IsUnaryFunction(_ T) (res bool) { return }
func (fn FunctionNumericComparative[T]) ProtoFunction() types.Function { return fn.FunctionNumeric }

func LT[T Numeric](v T) FunctionNumericComparative[T] {
	return FunctionNumericComparative[T]{ numeric(types.FunctionNumeric_LT, v) }
}
func GT[T Numeric](v T) FunctionComparative[T] {
	return FunctionNumericComparative[T]{ numeric(types.FunctionNumeric_GT, v) }
}

type FunctionBooleanCombinator[T any] struct { types.FunctionBoolean }

/*

type FunctionAddUint64 struct { types.FunctionNumeric }
var _ UnaryFunction[uint64, uint64] = FunctionAddUint64{}
func (FunctionAddUint64) IsUnaryFunction(_ uint64) (res uint64) { return }
func AddUint64(v uint64) FunctionAddUint64 { 
	return FunctionAddUint64{ types.FunctionNumeric { 
		Operator: types.FunctionNumeric_ADD,
		Operand: types.FunctionNumeric_Uint64Operand { Uint64Operand: v } 
	} } 
}

type FunctionAddInt64 struct { types.FunctionNumeric }
var _ UnaryFunction[int64, int64] = FunctionAddInt64{}
func (FunctionAddInt64) IsUnaryFunction(_ int64) (res int64) { return }
func AddInt64(v int64) FunctionAddInt64 { 
	return FunctionAddInt64{ types.FunctionNumeric { 
		Operator: typese.FunctionNumeric_ADD,
		Operand: types.FunctionAddNumeric_Int64Operand { Int64Operand: v } 
	} } 
}
*/

/*
type FunctionLTUint64 struct { types.FunctionNumeric }
var _ UnaryFunction[uint64, bool] = FunctionLTUint64{}
func (FuncitonLTUint64) IsUnaryFunction(_ uint64) (res bool) { return }
*/

type FunctionSelectField[T, U proto.Message] struct { types.FunctionSelectField }
var _ UnaryFucntion[T, U] = FunctionSelectField[T, U]{}
func (FunctionSelectField[T, U]) IsUnaryFunction(_ T) (res U) { return }
func (fn FunctionSelectField[T, U]) ProtoFunction() types.Function { return fn.FunctionSelectField }
func SelectField[T, U proto.Message](getChild func(T) U, number int32) FunctionSelectField[T, U] {
	field := T{}.ProtoReflect().Descriptor().Fields().ByNumber(number)
	return FunctionSelectField[T, U]{
		FieldNumber: field.Number(),
		FieldType: field.Kind(),
	}
}


// Sequence :: (T -> U) -> (U -> V) -> (T -> V)
type FunctionSequence[T, V any] struct { types.FunctionSequence }
var _ UnaryFunction[T, V] = FunctionSequence[T, V]{}
func (FunctionSequence[T, V]) IsUnaryFunction(_ T) (res V) { return }
func (fn FunctionSelectField[T, V]) ProtoFunction() types.Function { return fn.FunctionSequence }
func Seq[T, U, V any](f UnaryFunction[T, U], g UnaryFunction[U, V]) FunctionSequence[T, V] {
	// TODO: if f or g is FunctionSequence, append to it
	return FunctionSequence[T, V]{ types.FunctionSequence{ 
		Fn: []*types.Function{f.ProtoFunction(), g.ProtoFunction()};
	} }
}

// Filter :: (T -> bool) -> [T] -> [T]
type FunctionFilter[]
func Filter[T any](fn UnaryFunction[T, bool]) func(QueryBuilder[T]) QueryBuilder[T] {
	return func(input QueryBuilder[T]) QueryBuilder[T] {
		return queryBuild[T, T](input, types.QueryOperatorFilter{fn})
	}
}

// Map :: (T -> U) -> [T] -> [U]
func Map[T, U any](fn UnaryFunction[T, U]) func(QueryBuilder[T]) QueryBuilder[U] {
	return func(input QueryBuilder[T]) QueryBuilder[U] {
		return queryBuild[T, U](input, types.QueryOperatorMap{fn})
	}
}

// Fold :: ([T] -> U) -> [[T]] ->[U]
func Fold[T, U any](fn UnaryFunction[[]T, U]) func(QueryBuilder[[]T]) QueryBuilder[U] {
	return func(input QueryBuilder[[]T]) QueryBuilder[U] {
		return queryBuild[[]T, U](input, types.QueryOperatorFold{fn})
	}
}

// QueryBuilder T represents the Query that returns an enumarable of type T
type QueryBuilder[T any] struct { ops []types.FunctionEnumerable }

// queryBuild takes one query input and appends op, transformes it to another type of Query
func queryBuild[T, U any](input QueryBuilder[T], op types.QueryOperator) QueryBuilder[U] {
	input.ops = append(input.ops, op)
	return input
}

// Goal:
// query := func(xs Enumerable) Enumerable {
//   filtered := Filter(Seq(SelectField(MyDataType.GetSomeIntField, 3), GT(uint64(500))))(xs)
//   first100 := TakeFirst(100)(filtered)
//   return first100
// }