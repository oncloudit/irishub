package rpc

import (
	amino "github.com/tendermint/go-amino"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var cdc = amino.New()

func init() {
	ctypes.RegisterAmino(cdc)
}
