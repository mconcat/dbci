package dbci

import (
	dbm "github.com/tendermint/tm-db"
)

type KVStore interface {
	KVStoreRead
	KVStoreWrite
}

type Iterator interface {
	dbm.Iterator // TODO
}

type KVStoreRead interface {
	Get([]byte) []byte
	Iterator([]byte, []byte) Iterator
	ReverseIterator([]byte, []byte) Iterator
}

type DatabaseRead interface {
}

type KVStoreWrite interface {
	Set([]byte, []byte)
	Delete([]byte)
}

type DatabaseWrite interface {
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
