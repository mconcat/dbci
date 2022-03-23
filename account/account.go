package dbci

import (
	"fmt"
	"strings"
	"time"

	dbm "github.com/tendermint/tm-db"
)

type KeySegment interface {
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

type Key []KeySegment

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

type Account interface {
	ID() []byte // globally unique identifier for this account

	KVStore() dbm.KVStore
}

type AccountLocalUserStateAuth struct {
	// Copied from auth.pb.go/BaseAccount
	// should be moved
	Address State[string]
	PubKey State[[]byte]
	AccountNumber State[uint64]
	Sequence State[uint64]
}

type AccountLocalUserBank struct {
	// Copied from bank
	// should be moved
	AmountByDenom Mapping[sdk.Coin]

}

type AccountLocalUser struct {
	store dbm.KVStore
}

type AccountModule struct {
	store dbm.KVStore
}



// type AccountWASM struct

// type AccountInterchain struct
