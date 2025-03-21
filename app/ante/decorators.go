package ante

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	coinswaptypes "mods.irisnet.org/modules/coinswap/types"
	servicetypes "mods.irisnet.org/modules/service/types"
	tokenkeeper "mods.irisnet.org/modules/token/keeper"
	tokentypesv1 "mods.irisnet.org/modules/token/types/v1"
	tokentypesv1beta1 "mods.irisnet.org/modules/token/types/v1beta1"
)

// ValidateTokenDecorator is responsible for restricting the token participation of the swap prefix
type ValidateTokenDecorator struct {
	tk tokenkeeper.Keeper
}

// NewValidateTokenDecorator returns an instance of ValidateTokenDecorator
func NewValidateTokenDecorator(tk tokenkeeper.Keeper) ValidateTokenDecorator {
	return ValidateTokenDecorator{
		tk: tk,
	}
}

// AnteHandle checks the transaction
func (vtd ValidateTokenDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *ibctransfertypes.MsgTransfer:
			if containSwapCoin(msg.Token) {
				return ctx, sdkerrors.Wrap(errortypes.ErrInvalidRequest, "can't transfer coinswap liquidity tokens through the IBC module")
			}
		case *tokentypesv1.MsgBurnToken:
			if _, err := vtd.tk.GetToken(ctx, msg.Coin.Denom); err != nil {
				return ctx, sdkerrors.Wrap(errortypes.ErrInvalidRequest, "burnt failed, only native tokens can be burnt")
			}
		case *tokentypesv1beta1.MsgBurnToken:
			if _, err := vtd.tk.GetToken(ctx, msg.Symbol); err != nil {
				return ctx, sdkerrors.Wrap(errortypes.ErrInvalidRequest, "burnt failed, only native tokens can be burnt")
			}
		case *govv1.MsgSubmitProposal:
			if containSwapCoin(msg.InitialDeposit...) {
				return ctx, sdkerrors.Wrap(errortypes.ErrInvalidRequest, "can't deposit coinswap liquidity token for proposal")
			}
		case *govv1.MsgDeposit:
			if containSwapCoin(msg.Amount...) {
				return ctx, sdkerrors.Wrap(errortypes.ErrInvalidRequest, "can't deposit coinswap liquidity token for proposal")
			}
		}
	}
	return next(ctx, tx, simulate)
}

// ValidateServiceDecorator is responsible for checking the permission to execute MsgCallService
type ValidateServiceDecorator struct {
	SimulateTest bool
}

// NewValidateServiceDecorator returns an instance of ServiceAuthDecorator
func NewValidateServiceDecorator(simulateTest bool) ValidateServiceDecorator {
	return ValidateServiceDecorator{
		SimulateTest: simulateTest,
	}
}

// AnteHandle checks the transaction
func (vsd ValidateServiceDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	if vsd.SimulateTest {
		return next(ctx, tx, simulate)
	}

	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *servicetypes.MsgCallService:
			if msg.Repeated {
				return ctx, sdkerrors.Wrap(errortypes.ErrInvalidRequest, "currently does not support to create repeatable service invocation")
			}
		}
	}

	return next(ctx, tx, simulate)
}

func containSwapCoin(coins ...sdk.Coin) bool {
	for _, coin := range coins {
		if strings.HasPrefix(coin.Denom, coinswaptypes.LptTokenPrefix) {
			return true
		}
	}
	return false
}
