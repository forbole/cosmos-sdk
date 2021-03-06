package crypto

import (
	"github.com/tendermint/go-amino"
	tcrypto "github.com/tendermint/tendermint/crypto"
)

var cdc = amino.NewCodec()

func init() {
	RegisterAmino(cdc)
	tcrypto.RegisterAmino(cdc)
}

// RegisterAmino registers all go-crypto related types in the given (amino) codec.
func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterConcrete(PrivKeyLedgerSecp256k1{},
		"tendermint/PrivKeyLedgerSecp256k1", nil)
}
