package types 

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

func (fn FunctionNumeric) Execute() Numeric {
	switch operand := fn.Operand.(type) {
	case *FunctionNumeric_Uint64Operand:
		var x uint64
		err := proto.Unmarshal(v, &x)
		if err != nil {
			panic(err)
		}
		return ExecutePrimitiveNumeric(fn.Operator, operand.Uint64Operand, x)
	case *FunctionNumeric_Int64Operand:
		var x int64
		err := proto.Unmarshal(v, &x)
		if err != nil {
			panic(err)
		}
		return ExecutePrimitiveNumeric(fn.Operator, operand.Int64Operand, x)
	case *FunctionNumeric_Uint256Operand:
		var x Uint256
		err := proto.Unmarshal(v, &x)
		if err != nil {
			panic(err)
		}
		return ExecuteNonPrimitiveNumeric(fn.Operator, operand.Uint256Operand, x)
	case *FunctionNumeric_Int256Operand:
		var x Int256
		err := proto.Unmarshal(v, &x)
		if err != nil {
			panic(err)
		}
		return ExecuteNonPrimitiveNumeric(fn.Operator, operand.Int256Operand, x)
	}
}

func (fn FunctionSingle) Execute(v BytesSingle) BytesSingle {
	switch fn := fn.Function.(type) {
	case FunctionSingle_Numeric:
		return fn.Numeric.Execute(v)
	case FunctionSingle_Boolean:
		return fn.Boolean.Execute(v)
	case FunctionSingle_SelectField:
		return fn.SelectField.Execute(v)
	}
}

/*
func ExecuteSelectField(fn FunctionSelectField) (pre func(Single) Single, post func(FunctionEffectful) Single) {
	pre = func(v Single)
}
*/

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


func (from FromKey) Execute(store KVStore) BytesSingle {
	return MakeBytesSingle(from.Key, store.Get(from.Key))
}

func (op SingleFrom) Execute(store KVStore) BytesSingle {
	switch from := op.From.(type) {
	case SingleFrom_Key:
		return from.Key.Execute(store)
	}
}

func (op OperationSingle) ExecuteSelect(store KVStore, v BytesSingle) ([]byte, []byte) {
	return v.UnmodifiedRoot()
}

func (op OperationSingle) ExecuteSet(store KVStore, v BytesSingle) ([]byte, []byte) {
	key, value := v.ModifiedRoot()
	store.Set(key, value)
	return key, value
}

func (op OperationSingle) ExecuteDelete(store KVStore, v BytesSingle) []byte {
	key := v.Key()
	store.Delete(key)
	return key
}

func (op OperationSingle) Execute(accstore AccountStore) ([]byte, []byte) {
	store := accstore.Store(op.Account)

	v := op.From.Execute(store)
	for _, fn := range op.Functions {
		v = fn.Execute(v)
	}

	switch op.Type {
	case OperationSingle_SELECT:
		return op.ExecuteSelect(store. v)
	case OperationSingle_SET:
		return op.ExecuteSet(store, v)
	case OperationSingle_DELETE:
		return op.ExecuteDelete(store, v), nil
	}
}


type AccountStore interface {
	Store(account string) KVStore
}