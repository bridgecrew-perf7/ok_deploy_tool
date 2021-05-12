/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"poly-bridge/basedef"
	"poly-bridge/chainsdk"
	xecdsa "poly-bridge/utils/ecdsa"
	"poly-bridge/utils/files"
	"poly-bridge/utils/leveldb"
	"poly-bridge/utils/math"
	"poly-bridge/utils/wallet"
	"runtime"

	log "github.com/astaxie/beego/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/polynetwork/poly/native/service/header_sync/bsc"
	polyutils "github.com/polynetwork/poly/native/service/utils"
	oksdk "github.com/okex/exchain-go-sdk"
	"github.com/urfave/cli"
)

var (
	cfgPath  string
	cfg      = new(Config)
	cc       *ChainConfig
	storage  *leveldb.LevelDBImpl
	sdk      *chainsdk.EthereumSdk
	okclient oksdk.Client
	adm      *ecdsa.PrivateKey
	keystore string
)

const defaultAccPwd = "111111"

func setupApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "poly nftbridge deploy tool"
	app.Version = "1.0.0"
	app.Copyright = "Copyright in 2020 The Ontology Authors"
	app.Flags = []cli.Flag{
		LogLevelFlag,
		//LogDirFlag,
		ConfigPathFlag,
		ChainIDFlag,
		NFTNameFlag,
		NFTSymbolFlag,
		DstChainFlag,
		AssetFlag,
		DstAssetFlag,
		SrcAccountFlag,
		DstAccountFlag,
		//FeeTokenFlag,
		//ERC20TokenFlag,
		//NativeTokenFlag,
		AmountFlag,
		TokenIdFlag,
		MethodCodeFlag,
		OwnerAccountFlag,
		AdminIndexFlag,
		AddGasFlag,
		EpochFlag,
	}
	app.Commands = []cli.Command{
		CmdSample,
		CmdDeployECCDContract,
		CmdDeployECCMContract,
		CmdDeployCCMPContract,
		CmdBindERC20Asset,
		CmdTransferECCDOwnership,
		CmdTransferECCMOwnership,
		CmdRegisterSideChain,
		CmdApproveSideChain,
		CmdSyncSideChainGenesis2Poly,
		CmdSyncPolyGenesis2SideChain,
		CmdNativeBalance,
		CmdNativeTransfer,
		CmdEnv,
	}

	app.Before = beforeCommands
	app.After = afterCommond
	return app
}

func main() {

	app := setupApp()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// action execute after commands
func beforeCommands(ctx *cli.Context) (err error) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// load config instance
	cfgPath = ctx.GlobalString(getFlagName(ConfigPathFlag))
	if err = files.ReadJsonFile(cfgPath, cfg); err != nil {
		return fmt.Errorf("read config json file, err: %v", err)
	}

	//logDir := ctx.GlobalString(getFlagName(LogDirFlag))
	//logFormat := fmt.Sprintf(`{"filename":"%s/deploy.log", "perm": "0777"}`, logDir)
	loglevel := ctx.GlobalUint64(getFlagName(LogLevelFlag))
	logFormat := fmt.Sprintf(`{"level:":"%d"}`, loglevel)
	if err := log.SetLogger("console", logFormat); err != nil {
		return fmt.Errorf("set logger failed, err: %v", err)
	}

	// prepare storage for persist account passphrase
	storage = leveldb.NewLevelDBInstance(cfg.LevelDB)

	// select src chainID and prepare config and accounts
	chainID := ctx.GlobalUint64(getFlagName(ChainIDFlag))
	selectChainConfig(chainID)

	if _, err := os.Stat(cfg.Keystore); os.IsNotExist(err) {
		return fmt.Errorf("keystore dir %s is not exist", cfg.Keystore)
	}
	keystore = cfg.Keystore
	admIndex := ctx.GlobalInt(getFlagName(AdminIndexFlag))
	if admIndex < 0 || admIndex >= len(cfg.AdminAccountList) {
		return fmt.Errorf("admin index out of range")
	}
	admAddr := cfg.AdminAccountList[admIndex]
	if adm, err = wallet.LoadEthAccount(storage, keystore, admAddr, defaultAccPwd); err != nil {
		return fmt.Errorf("load eth account for chain %d faild, err: %v", cc.SideChainID, err)
	}

	uAddGas := ctx.GlobalUint64(getFlagName(AddGasFlag))
	if uAddGas > 0 {
		chainsdk.DefaultAddGasPrice = new(big.Int).SetUint64(uAddGas * 1000000000)
	}

	if chainID == basedef.OK_CROSSCHAIN_ID {
		okconfig, _ := oksdk.NewClientConfig(cc.RPC, "okexchain-65", oksdk.BroadcastBlock, "0.01okt", 200000, 0, "")
		okclient = oksdk.NewClient(okconfig)
	} else {
		if sdk, err = chainsdk.NewEthereumSdk(cc.RPC); err != nil {
			return fmt.Errorf("generate sdk for chain %d faild, err: %v", cc.SideChainID, err)
		}
	}

	return nil
}

func afterCommond(ctx *cli.Context) error {
	log.Info("\r\n" +
		"\r\n" +
		"\r\n")
	return nil
}

func handleSample(ctx *cli.Context) error {
	log.Info("start to debug sample...")
	//feeToken := ctx.BoolT(getFlagName(FeeTokenFlag))
	//nativeToken := ctx.Bool(getFlagName(NativeTokenFlag))
	//log.Info("feeToken %v, nativeToken %v", feeToken, nativeToken)
	//getFeeTokenOrERC20Asset(ctx)
	return nil
}

func handleCmdDeployECCDContract(ctx *cli.Context) error {
	log.Info("start to deploy eccd contract...")

	addr, err := sdk.DeployECCDContract(adm)
	if err != nil {
		return fmt.Errorf("deploy eccd for chain %d failed, err: %v", cc.SideChainID, err)
	}

	cc.ECCD = addr.Hex()
	log.Info("deploy eccd for chain %d success %s", cc.SideChainID, addr.Hex())
	return updateConfig()
}

func handleCmdDeployECCMContract(ctx *cli.Context) error {
	log.Info("start to deploy eccm contract...")

	eccd := common.HexToAddress(cc.ECCD)
	addr, err := sdk.DeployECCMContract(adm, eccd, cc.SideChainID)
	if err != nil {
		return fmt.Errorf("deploy eccm for chain %d failed, err: %v", cc.SideChainID, err)
	}
	cc.ECCM = addr.Hex()
	log.Info("deploy eccm for chain %d success %s", cc.SideChainID, addr.Hex())
	return updateConfig()
}

func handleCmdDeployCCMPContract(ctx *cli.Context) error {
	log.Info("start to deploy ccmp contract...")

	eccm := common.HexToAddress(cc.ECCM)
	addr, err := sdk.DeployECCMPContract(adm, eccm)
	if err != nil {
		return fmt.Errorf("deploy ccmp for chain %d failed, err: %v", cc.SideChainID, err)
	}
	cc.CCMP = addr.Hex()
	log.Info("deploy ccmp for chain %d success %s", cc.SideChainID, addr.Hex())
	return updateConfig()
}

func handleCmdBindERC20Asset(ctx *cli.Context) error {
	log.Info("start to bind nft asset...")

	srcAsset := flag2address(ctx, AssetFlag)
	dstAsset := flag2address(ctx, DstAssetFlag)
	dstChainId := flag2Uint64(ctx, DstChainFlag)
	dstChainCfg := customSelectChainConfig(dstChainId)
	owner := xecdsa.Key2address(adm)
	proxy := common.HexToAddress(cc.LockProxy)

	hash, err := sdk.BindERC20Asset(
		adm,
		proxy,
		srcAsset,
		dstAsset,
		dstChainId,
	)
	if err != nil {
		return fmt.Errorf("bind erc20 asset (src chain id %d, src asset %s, src proxy %s) - "+
			"(dst chain id %d, dst asset %s, dst proxy %s)"+
			" for user %s failed, err: %v",
			cc.SideChainID, srcAsset.Hex(), cc.LockProxy,
			dstChainId, dstAsset.Hex(), dstChainCfg.LockProxy,
			owner.Hex(), err)
	}

	log.Info("bind erc20 asset (src chain id %d, src asset %s, src proxy %s) - "+
		"(dst chain id %d, dst asset %s, dst proxy %s)"+
		" for user %s success! txhash %s",
		cc.SideChainID, srcAsset.Hex(), cc.LockProxy,
		dstChainId, dstAsset.Hex(), dstChainCfg.LockProxy,
		owner.Hex(), hash.Hex())
	return nil
}

func handleCmdTransferECCDOwnership(ctx *cli.Context) error {
	log.Info("start to transfer eccd ownership...")

	eccd := common.HexToAddress(cc.ECCD)
	eccm := common.HexToAddress(cc.ECCM)

	if hash, err := sdk.TransferECCDOwnership(adm, eccd, eccm); err != nil {
		return fmt.Errorf("transfer eccd %s ownership to eccm %s on chain %d failed, err: %v",
			cc.ECCD, cc.ECCM, cc.SideChainID, err)
	} else {
		log.Info("transfer eccd %s ownership to eccm %s on chain %d success, txhash: %s",
			cc.ECCD, cc.ECCM, cc.SideChainID, hash.Hex())
	}
	return nil
}

func handleCmdTransferECCMOwnership(ctx *cli.Context) error {
	log.Info("start to transfer eccm ownership...")

	eccm := common.HexToAddress(cc.ECCM)
	ccmp := common.HexToAddress(cc.CCMP)

	if hash, err := sdk.TransferECCMOwnership(adm, eccm, ccmp); err != nil {
		return fmt.Errorf("transfer eccm %s ownership to ccmp %s on chain %d failed, err: %v",
			cc.ECCM, cc.CCMP, cc.SideChainID, err)
	} else {
		log.Info("transfer eccm %s ownership to ccmp %s on chain %d success, txhash: %s",
			cc.ECCM, cc.CCMP, cc.SideChainID, hash.Hex())
	}

	return nil
}

func handleCmdRegisterSideChain(ctx *cli.Context) error {
	validators, err := wallet.LoadPolyAccountList(cfg.Poly.Keystore, cfg.Poly.Passphrase)
	if err != nil {
		return err
	}
	polySdk, err := chainsdk.NewPolySdkAndSetChainID(cfg.Poly.RPC)
	if err != nil {
		return err
	}

	// todo: 验证heco的注册方式
	eccd := common.HexToAddress(cc.ECCD)
	chainID := cc.SideChainID
	switch chainID {
	case basedef.ETHEREUM_CROSSCHAIN_ID:
		router := polyutils.ETH_ROUTER
		err = polySdk.RegisterSideChain(validators[0], chainID, 1, router, eccd, cc.SideChainName)

	case basedef.BSC_CROSSCHAIN_ID:
		router := polyutils.BSC_ROUTER
		ext := bsc.ExtraInfo{
			ChainID: new(big.Int).SetUint64(chainID),
		}
		extEnc, _ := json.Marshal(ext)
		err = polySdk.RegisterSideChainExt(validators[0], chainID, 1, router, eccd, cc.SideChainName, extEnc)

	case basedef.HECO_CROSSCHAIN_ID:
		router := polyutils.HECO_ROUTER
		err = polySdk.RegisterSideChain(validators[0], chainID, 1, router, eccd, cc.SideChainName)

	case basedef.OK_CROSSCHAIN_ID:
		router := uint64(12)
		err = polySdk.RegisterSideChain(validators[0], chainID, 1, router, eccd, cc.SideChainName)

	default:
		err = fmt.Errorf("chain id %d invalid", chainID)
	}

	if err != nil {
		return err
	}

	log.Info("register side chain %d eccd %s success", chainID, eccd.Hex())
	return nil
}

func handleCmdApproveSideChain(ctx *cli.Context) error {
	validators, err := wallet.LoadPolyAccountList(cfg.Poly.Keystore, cfg.Poly.Passphrase)
	if err != nil {
		return err
	}
	polySdk, err := chainsdk.NewPolySdkAndSetChainID(cfg.Poly.RPC)
	if err != nil {
		return err
	}
	if err := polySdk.ApproveRegisterSideChain(cc.SideChainID, validators); err != nil {
		return fmt.Errorf("failed to approve register side chain, err: %s", err)
	}

	log.Info("approve register side chain %d success", cc.SideChainID)
	return nil
}

func handleCmdSyncSideChainGenesis2Poly(ctx *cli.Context) error {
	log.Info("start to sync side chain %s genesis header to poly chain...", cc.SideChainName)

	polySdk, err := chainsdk.NewPolySdkAndSetChainID(cfg.Poly.RPC)
	if err != nil {
		return err
	}
	validators, err := wallet.LoadPolyAccountList(cfg.Poly.Keystore, cfg.Poly.Passphrase)
	if err != nil {
		return err
	}

	switch cc.SideChainID {
	case basedef.ETHEREUM_CROSSCHAIN_ID:
		err = SyncEthGenesisHeader2Poly(cc.SideChainID, sdk, polySdk, validators)
	case basedef.BSC_CROSSCHAIN_ID:
		err = SyncBscGenesisHeader2Poly(cc.SideChainID, sdk, polySdk, validators)
	case basedef.HECO_CROSSCHAIN_ID:
		err = SyncHecoGenesisHeader2Poly(cc.SideChainID, sdk, polySdk, validators)
	case basedef.OK_CROSSCHAIN_ID:
		epoch := flag2Uint64(ctx, EpochFlag)
		err = SyncOKGenesisHeader2Poly(cc.SideChainID, okclient, polySdk, validators, int64(epoch))
	default:
		err = fmt.Errorf("chain id %d invalid", cc.SideChainID)
	}
	if err != nil {
		return fmt.Errorf("sync side chain %d genesis header to poly failed, err: %v", cc.SideChainID, err)
	} else {
		log.Info("sync side chain %d genesis header to poly success!", cc.SideChainID)
	}
	return nil
}

func handleCmdSyncPolyGenesis2SideChain(ctx *cli.Context) error {
	log.Info("start to sync poly chain genesis header to side chain...")

	polySdk, err := chainsdk.NewPolySdkAndSetChainID(cfg.Poly.RPC)
	if err != nil {
		return err
	}
	eccm := common.HexToAddress(cc.ECCM)

	if err := SyncPolyGenesisHeader2Eth(
		polySdk,
		adm,
		sdk,
		eccm,
	); err != nil {
		return fmt.Errorf("sync poly chain genesis header to side chain %d failed, err: %v", cc.SideChainID, err)
	}
	log.Info("sync poly chain genesis header to side chain %d success!", cc.SideChainID)
	return nil
}

func handleGetNativeBalance(ctx *cli.Context) error {
	owner := flag2address(ctx, SrcAccountFlag)
	balance, err := sdk.GetNativeBalance(owner)
	if err != nil {
		return fmt.Errorf("get native balance faild, err: %v", err)
	}
	log.Info("%s native balance is %s", owner.Hex(), balance.String())
	return nil
}

func handleCmdNativeTransfer(ctx *cli.Context) error {
	log.Info("start to transfer native token on chain %s...", cc.SideChainName)

	from := flag2address(ctx, SrcAccountFlag)
	key, err := wallet.LoadEthAccount(storage, keystore, from.Hex(), defaultAccPwd)
	if err != nil {
		return err
	}

	to := flag2address(ctx, DstAccountFlag)
	amount := flag2big(ctx, AmountFlag)
	tx, err := sdk.TransferNative(key, to, amount)
	if err != nil {
		return err
	}
	log.Info("%s transfer %s to %s success, txhash %s", from.Hex(), amount.String(), to.Hex(), tx.Hex())
	return nil
}

func handleCmdEnv(ctx *cli.Context) error {
	currentInfo := fmt.Sprintf("current env: side chain name %s, side chain id %d\r\n", cc.SideChainName, cc.SideChainID)

	polyInfo := fmt.Sprintf("poly side chain id - %d\r\n", basedef.POLY_CROSSCHAIN_ID)
	ethInfo := fmt.Sprintf("eth side chain id - %d\r\n", basedef.ETHEREUM_CROSSCHAIN_ID)
	ontInfo := fmt.Sprintf("ont side chain id - %d\r\n", basedef.ONT_CROSSCHAIN_ID)
	neoInfo := fmt.Sprintf("neo side chain id - %d\r\n", basedef.NEO_CROSSCHAIN_ID)
	bscInfo := fmt.Sprintf("bsc side chain id - %d\r\n", basedef.BSC_CROSSCHAIN_ID)
	hecoInfo := fmt.Sprintf("heco side chain id - %d\r\n", basedef.HECO_CROSSCHAIN_ID)
	o3Info := fmt.Sprintf("o3 side chain id - %d\r\n", basedef.O3_CROSSCHAIN_ID)

	log.Info(currentInfo, polyInfo, ethInfo, ontInfo, neoInfo, bscInfo, hecoInfo, o3Info)

	owner := flag2address(ctx, OwnerAccountFlag)
	addr := owner.Hex()
	log.Info("check your owner address %s in dir %s", keystore, addr)
	_, err := wallet.LoadEthAccount(storage, keystore, addr, defaultAccPwd)
	if err != nil {
		return err
	}
	//enc := crypto.FromECDSA(key)
	//log.Info(hex.EncodeToString(enc))
	return nil
}

// getFeeTokenOrERC20Asset return feeToken if `feeToken` is true
//func getFeeTokenOrERC20Asset(ctx *cli.Context) common.Address {
//	if ctx.Bool(getFlagName(NativeTokenFlag)) {
//		return nativeToken
//	}
//	if ctx.Bool(getFlagName(ERC20TokenFlag)) {
//		return flag2address(ctx, ERC20TokenFlag)
//	}
//	return common.HexToAddress(cc.FeeToken)
//}

func updateConfig() error {
	if err := files.WriteJsonFile(cfgPath, cfg, true); err != nil {
		return err
	}
	log.Info("update config success!", cfgPath)
	return nil
}

func selectChainConfig(chainID uint64) {
	cc = customSelectChainConfig(chainID)
}

func flag2string(ctx *cli.Context, f cli.Flag) string {
	fn := getFlagName(f)
	data := ctx.String(fn)
	return data
}

func flag2address(ctx *cli.Context, f cli.Flag) common.Address {
	data := flag2string(ctx, f)
	return common.HexToAddress(data)
}

func flag2big(ctx *cli.Context, f cli.Flag) *big.Int {
	fn := getFlagName(f)
	data := ctx.String(fn)
	return math.String2BigInt(data)
}

func flag2Uint64(ctx *cli.Context, f cli.Flag) uint64 {
	fn := getFlagName(f)
	data := ctx.Uint64(fn)
	return data
}

func customSelectChainConfig(chainID uint64) *ChainConfig {
	switch chainID {
	case basedef.ETHEREUM_CROSSCHAIN_ID:
		return cfg.Ethereum
	case basedef.BSC_CROSSCHAIN_ID:
		return cfg.Bsc
	case basedef.HECO_CROSSCHAIN_ID:
		return cfg.Heco
	case basedef.OK_CROSSCHAIN_ID:
		return cfg.Ok
	}
	panic(fmt.Sprintf("invalid chain id %d", chainID))
}
