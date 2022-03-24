package account_test

import (
	"testing"

	"github.com/mconcat/dbci/account"
)

type State[T any] interface {
	Get() T	
	Set(T)
}

/*
type AccountStateAuth struct {
	// Copied from auth.pb.go/BaseAccount
	// should be moved
	Address State[string]
	PubKey State[[]byte]
	AccountNumber State[uint64]
	Sequence State[uint64]
}

func NewAccountStateAuth()  {

}
*/

type AccountStateBank struct {
	// Copied from bank
	// should be moved
	AmountByDenom Mapping[sdk.Coin]
}

func NewAccountStateBank(acc account.Account) AccountStateBank {
	return AccountStateBank{
		AmountByDenom: NewMapping[sdk.Coin](acc.KVStore(), "bank", "amount"),
	}
}

func (acc AccountStateBank) GetBalance(denom string) State[sdk.Coin] {
	return acc.AmountByDenom.Of(denom)
}

func (acc AccountStateBank) Send(acc2 AccountStateBank, coin sdk.Coin) {
	return acc.AmountByDenom.Of(coin.Denom).MoveTo(acc2.AmountByDenom.Of(coin.Denom), coin.Amount)
/*
	Translated into:

	Atomic {
		Set {
			FromSingle acc.ID + bank + amount + coin.Denom :: Uint256
			Sub coin.Amount :: Uint256
		}
		Set {
			FromSingle acc2.ID + bank + amount + coin.Denom :: Uint256
			Add coin.Amount :: Uint256
		}
	}
*/
}