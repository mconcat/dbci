package types_test

import (
	"github.com/mconcat/dbci/operation/types"
)

type CaseFunctionNumeric struct {
	Function types.FunctionNumeric
	ExpectedBehaviour func(x Numeric) Numeric
}

func casesFunctionNumeric() []CaseFunctionNumeric {
	return []CaseFunctionNumeric{{ 
		types.FunctionNumeric{ 
			Operator: types.FunctionNumeric_ADD,
			Operand: &types.FunctionNumeric_Uint256Operand{&types.Uint256{/*TODO*/}},
		},
		func(x Numeric) Numeric {
			return x.(Uint256).Add(/*TODO*/)
		}, 
	}, {
		types.FunctionNumeric{
			Operator: types.FunctionNumeric_SUB,
			Operand: &types.FunctionNumeric_Uint256Operand{&types.Uint256{/*TODO*/}},
		},
		func(x Numeric) Numeric {
			return x.(Uint256).SUB(/*TODO*/)
		},
	}}
}

type CaseFromKey struct {
	From types.FromKey
	ExpectedBehaviour func(store types.KVStore) ([]byte, []byte)
}

type CaseOperationSingle struct {
	Operation types.OperationSingle
	ExpectedBehaviour func(store types.KVStore) bool
}