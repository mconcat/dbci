package proto

import "github.com/mconcat/dbci/query"

func (q Query) handle(op query.Operator) interface{} {
	switch op.(type) {
	case QueryOperatorMap[T, U]
	}
}