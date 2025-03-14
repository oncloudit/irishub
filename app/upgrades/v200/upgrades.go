package v200

import (
	"context"
	"fmt"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/upgrade/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"

	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/evmos/ethermint/x/feemarket"
	feemarkettypes "github.com/evmos/ethermint/x/feemarket/types"

	"github.com/irisnet/irishub/v4/app/upgrades"
	irisevm "github.com/irisnet/irishub/v4/modules/evm"
)

// Upgrade defines a struct containing necessary fields that a SoftwareUpgradeProposal
var Upgrade = upgrades.Upgrade{
	UpgradeName:               "v2.0",
	UpgradeHandlerConstructor: upgradeHandlerConstructor,
	StoreUpgrades: &storetypes.StoreUpgrades{
		Added:   []string{evmtypes.StoreKey, feemarkettypes.StoreKey},
		Deleted: []string{icahosttypes.StoreKey},
	},
}

func upgradeHandlerConstructor(
	m *module.Manager,
	c module.Configurator,
	box upgrades.Toolbox,
) upgradetypes.UpgradeHandler {
	return func(context context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(context)
		fromVM[evmtypes.ModuleName] = irisevm.AppModule{}.ConsensusVersion()
		fromVM[feemarkettypes.ModuleName] = feemarket.AppModule{}.ConsensusVersion()

		if err := box.EvmKeeper.SetParams(ctx, evmParams); err != nil {
			return nil, err
		}

		if err := box.FeeMarketKeeper.SetParams(ctx, generateFeemarketParams(ctx.BlockHeight())); err != nil {
			return nil, err
		}

		//transfer token ownership
		owner, err := sdk.AccAddressFromBech32(evmToken.Owner)
		if err != nil {
			return nil, err
		}
		if err := box.TokenKeeper.UnsafeTransferTokenOwner(ctx, evmToken.Symbol, owner); err != nil {
			return nil, err
		}

		//update consensusParams.Block.MaxGas
		consensusParams := box.ReaderWriter.GetConsensusParams(ctx)
		consensusParams.Block.MaxGas = maxBlockGas
		box.ReaderWriter.StoreConsensusParams(ctx, consensusParams)

		//add Burner Permission for authtypes.FeeCollectorName
		feeModuleAccount := box.AccountKeeper.GetModuleAccount(ctx, authtypes.FeeCollectorName)
		account, ok := feeModuleAccount.(*authtypes.ModuleAccount)
		if !ok {
			return nil, fmt.Errorf("feeCollector accountis not *authtypes.ModuleAccount")
		}
		account.Permissions = append(account.Permissions, authtypes.Burner)
		box.AccountKeeper.SetModuleAccount(ctx, account)

		// delete ica moudule version from upgrade moudule
		store := ctx.KVStore(box.GetKey(upgradetypes.StoreKey))
		versionStore := prefix.NewStore(store, []byte{types.VersionMapByte})
		versionStore.Delete([]byte(icatypes.ModuleName))

		return box.ModuleManager.RunMigrations(ctx, c, fromVM)
	}
}
