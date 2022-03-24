package types 

import (
	"github.com/mconcat/dbci/operation/types"
)

func ExecutePrimitiveNumeric[T PrimitiveNumeric](operator FunctionNumeric_Operator, operand, v T) T {
	switch operator {
	case FunctionNumeric_ADD:
		return v+operand
	case FunctionNumeric_SUB:
		return v-operand
	case FunctionNumeric_SUBFROM:
		return operand-v
	}
}

func ExecuteNonPrimitiveNumeric[T interface{ ~Int256 | ~Uint256 }](operator FunctionNumeric_Operator, operand, v T) T {
	switch operator {
	case FunctionNumeric_ADD:
		return v.Add(operand)
	case FunctionNumeric_SUB:
		return v.Sub(operand)
	case FunctionNumeric_SUBFROM:
		return operand.Sub(v)
	}
}

func ExecuteNumeric(fn FunctionNumeric, v Numeric) Numeric {
	switch operand := fn.Operand.(type) {
	case *FunctionNumeric_Uint64Operand:
		return ExecutePrimitiveNumeric(fn.Operator, operand.Uint64Operand, v.(uint64))
	case *FunctionNumeric_Int64Operand:
		return ExecutePrimitiveNumeric(fn.Operator, operand.Int64Operand, v.(int64))
	case *FunctionNumeric_Uint256Operand:
		return ExecuteNonPrimitiveNumeric(fn.Operator, operand.Uint256Operand, v.(Uint256))
	case *FunctionNumeric_Int256Operand:
		return ExecuteNonPrimitiveNumeric(fn.Operator, operand.Int256Operand, v.(Uint256))
	}
}

func ExecuteSelect(v Single, fns []Function) Single {
	
}

// Any of the operations are lazy evaluated, at the time when the actual database manipulation need to happen(Select, Set and Delete)
// Effectful function happens on 
func ExecuteEffectful(fn FunctionEffectful, v Single, fns []Function) Single {
	switch fn.FunctionType {
	case FunctionEffectful_SELECT:
		return ExecuteSelect(v, fns)
	case FunctionEffectful_SET:
		return ExecuteSet(v, fns)
	case FunctionEffectful_DELETE:
		return ExecuteDelete(v, fns)
	}
}

func ExecuteSelectField(fn FunctionSelectField) (pre func(Single) Single, post func(FunctionEffectful) Single) {
	pre = func(v Single)
}

func ExecuteEnumerable(fn FunctionEnumerable, v RawEnumerator) RawEnumerator {
	switch fn.FunctionType {
	case FunctionEnumerable_MAP:
		return MakeRawMapEnumerator(fn.Fn, v)
	case FunctionEnumerable_FOREACH:
		return MakeRawForeachEnumerator(fn.Fn, v)
	case FunctionEnumerable_TAKEFIRST:
		return MakeRawTakefirstSingle(v)
	}
}
