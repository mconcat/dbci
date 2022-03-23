package operation

import (
	"github.com/mconcat/dbci/operation/types"
)

type Enumarable interface {
	Next() ([]byte, []byte)
}

func ExecutePrimitiveNumeric[T PrimitiveNumeric](operator types.FunctionNumeric_Operator, operand, v T) T {
	switch operator {
	case types.FunctionNumeric_ADD:
		return v+operand
	case types.FunctionNumeric_SUB:
		return v-operand
	}
}

func ExecuteNumeric(fn types.FunctionNumeric, v Numeric) interface{} {
	switch operand := fn.Operand.(type) {
	case *types.FunctionNumeric_Uint64Operand:
		return ExecutePrimitiveNumeric(fn.Operator, operand.Uint64Operand, v.(uint64))
	case *types.FunctionNumeric_Int64Operand:
		return ExecutePrimitiveNumeric(fn.Operator, operand.Int64Operand, v.(int64))
	}
}
