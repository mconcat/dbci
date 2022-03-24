package types 

import (
	dbm "github.com/tendermint/tm-db"
)

type Index interface {
	KeyString() string
}

type Uint64 uint64

func (key Uint64) KeyString() string {
	return fmt.Sprintf("%016x", uint64(key))
}

type Bytes []byte

func (key Bytes) KeyString() string {
	return fmt.Sprintf("%x", []byte(key))
}

type Time time.Time

type String string

func (key String) KeyString() string {
	return string(key)
}

type Key []Index

func (key Key) String() string {
	keys := make([]string, len(key))
	for i, k := range key {
		keys[i] = k.KeyString()
	}
	return strings.Join(keys, "/")
}

func (key Key) Bytes() []byte {
	return []byte(key.String())
}

func (key Key) KeyString() string { return key.String() }

// C# LINQ Enumerable<T>
// ReactiveX Obervable<T>
type Enumerator[K Index, V any] interface {
	Peek() (K, V)
	Next()
	IsEmpty() bool
}

type BytesEnumerator = Enumerator[Index, []byte]

type SliceEnumerator[K Index, V any] struct {
	Elems []Pair[K, V]
}

func (enum *SliceEnumerator[K, V]) Peek() (K, V) {
	res := enum.Elems[0]
	return res.Key, res.Value
}

func (enum *SliceEnumerator[K, V]) Next() {
	enum.Elems = enum.Elems[1:]
}

func (enum *SliceEnumerator[K, V]) IsEmpty() bool {
	return len(enum.Elems) == 0
}

// MergedEnumerator merges enumerators by key order
type MergedEnumerator[K Index, V any] struct {
	Children []Enumerator[K, V]

	cachedChild *int
}

func (enum *MergedEnumerator[K, V]) Peek() (K, V) {
	if enum.cachedChild == nil {
		var minChild int
		var minChildKey []byte
		for i, child := range enum.Children {
			if child.IsEmpty() {
				continue
			}
			childKey := child.Peek().Key.KeyString()
			if bytes.Compare([]byte(childKey), []byte(minChildKey)) == -1 {
				minChild = i
				minChildKey = childKey
			}
		}	
		enum.cachedChild = &minChild
	}

	res := enum.Children[*enum.cachedChild].Peek()
	return res.Key, res.Value
}

func (enum *MergedEnumerator[K, V]) Next() {
	enum.Children[*enum.cachedChild].Next()
}

func (enum *MergedEnumerator[K, V]) IsEmpty() bool {
	for i, child := range enum.Children {
		if !child.IsEmpty() {
			return false
		}
	}
	return true
}

// ConcatEnumerator sequentially returns from multiple enumerators
type ConcatEnumerator[K Index, V any] struct {
	Children []Enumerator[K, V]	
}

// MapEnumerator applies operationtypes.Function to each element
type MapEnumerator[K Index, V any] struct {
	Child Enumerator[K, V]
	Function Function
}

func MakeRawMapEnumerator(fn Function, enum RawEnumerator) RawEnumerator {
	return MapEnumerator[[]byte, []byte] {
		Child: enum,
		Funciton: fn,
	}
}

func (enum *MapEnumerator[K, V]) Peek() (K, V) {
	peek := enum.Child.Peek()
	return Pair{peek.Key, ExecuteFunction(peek.Value).(V)}
}

func (enum *MapEnumerator[K, V]) Next() (K, V) {
	enum.Child.Next()
}

func (enum *MapEnumerator[K, V]) IsEmpty() bool {
	return enum.Child.IsEmpty()
}

// K is the type of key
// R is the type of root value, stored under the key in the state
// V is the type of the current value, possibly manipulated and field-selected
type Single[K Index, R any, V any] interface {
	UnmodifiedRoot() (K, R) // returns the original root key-value pair
	ModifiedRoot() (K, R) // returns the root key-value pair as manipulated
	Key() K // return the key to be deleted 
	Value() V
}

type BytesSingle = Single[[]byte, []byte, []byte]

func MakeBytesSingle(key []byte, value []byte) BytesSingle {
	return Pair{
		Key: key,
		Value: value,
	}
}

// ReactiveX Single<T>
type Pair[K Index, V any] struct {
	Key K
	Value V
}

func (pair Pair[K, V]) Key() K {
	return pair.Key
} 

func (pair Pair[K, V]) Value() V {
	return pair.Value
}

func (pair Pair[K, V]) UnmodifiedRoot() (K, V) {
	return pair.Key, pair.Value
}

func (pair Pair[K, V]) ModifiedRoot() (K, V) {
	return pair.Key, pair.Value
}

// FieldSelectedPair represents a pointer to a field in a protobuf message.
// When Set(), FieldSelectedPair replaces the field value of the root pair
// to the current value. This is done by ConsumeTag() function of protobuf wire utility.
type FieldSelectedPair[K Index, V any, F any] struct {
	Pair[K, F]
	FieldSelector FunctionSelectField
	Root V
}

func (pair *FieldSelectedPair[K, V, F]) UnmodifiedRoot() (K, V) {
	return pair.Key, pair.Root
}

func (pair *FieldSelectedPair[K, V, F]) ModifiedRoot() (K, V) {
	// TODO: set field and return
}

func (pair *FieldSelectedPair[K, V, F]) Key() K {
	return pair.Key
}

func (pair *FieldSelectedPair[K, V, F]) Value() F {
	return pair.Value
}

type KVStore interface {
	dbm.KVStore // TODO
}

type Iterator interface {
	dbm.Iterator // TODO
}

// Used for onchain data. Queryable. Part of consensus.
type KVStoreLocal struct {
}

// Used for read-locked onchain data. Queryable. Part of consensus. Get only.
type KVStoreReadonly struct {
}

// Used for offchain data. Not queryable. Not a part of consensus. Set only.
type KVStoreOffchain struct {
}

// Concurrency control
type KVStoreConcurrent struct {

}