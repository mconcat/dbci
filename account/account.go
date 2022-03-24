package dbci

import (
	"fmt"
	"strings"
	"time"

	dbm "github.com/tendermint/tm-db"
)


type Account interface {
	ID() []byte // globally unique identifier for this account

	KVStore() dbm.KVStore
}

type AccountLocalUser struct {
	store dbm.KVStore
}

type AccountModule struct {
	store dbm.KVStore
}

// type AccountWASM struct

// type AccountInterchain struct

